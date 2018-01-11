package gghelper

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/greengrass"
	"github.com/aws/aws-sdk-go/service/lambda"
)

// LambdaConfig - configuration for creating or updating Lambda functions
type LambdaConfig struct {
	Directory string
	Handler   string
	Name      string
	Pinned    bool
	RoleArn   string
	Runtime   string
}

// zip - zip a directory for lambda
func zipDirectory(sourceDir string) (*bytes.Buffer, error) {
	// Make sure the source directory exists and is a directory
	info, err := os.Stat(sourceDir)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		// XXX
		return nil, nil
	}

	// Create a memory backed zip file for the contents
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Print(err)
			return nil
		}

		// The zip archive nees to contain relative filenames
		rel, _ := filepath.Rel(sourceDir, path)
		if rel == "." {
			return nil
		}

		header, err := zip.FileInfoHeader(info)

		// Directories get a header entry with a trailing "/"
		if info.IsDir() {
			// fmt.Printf("Adding directory %s as %s\n", path, rel)
			header.Name = rel + "/"
			_, err := w.CreateHeader(header)
			return err
		}

		// Add the file to the zip archive
		// fmt.Printf("Adding file %s as %s\n", path, rel)
		header.Name = rel
		header.Method = zip.Deflate
		writer, err := w.CreateHeader(header)
		f, err := os.Open(path)
		defer f.Close()
		_, err = io.Copy(writer, f)
		return err
	})

	// Done creating the zip file
	w.Close()

	// Write it out as a test
	f, err := os.Create("foo.zip")
	f.Write(buf.Bytes())
	f.Close()
	return buf, nil
}

// SubmitFunction - xxx
func (ggSession *GreengrassSession) SubmitFunction(lambdaConfig LambdaConfig) error {
	var functionArn string
	latest := "latest"

	name := lambdaConfig.Name
	functionInput := lambda.GetFunctionInput{
		FunctionName: &name,
	}

	createFunction := true
	_, err := ggSession.lambda.GetFunction(&functionInput)
	if err == nil {
		createFunction = false
	}

	zip, err := zipDirectory(lambdaConfig.Directory)
	if err != nil {
		return err
	}

	if createFunction {
		fmt.Printf("Creating function: %s\n", lambdaConfig.Name)
		functionCode := lambda.FunctionCode{
			ZipFile: zip.Bytes(),
		}
		createFunctionInput := lambda.CreateFunctionInput{
			FunctionName: &lambdaConfig.Name,
			Runtime:      &lambdaConfig.Runtime,
			Handler:      &lambdaConfig.Handler,
			Role:         &lambdaConfig.RoleArn,
			Code:         &functionCode,
		}
		_, err := ggSession.lambda.CreateFunction(&createFunctionInput)
		if err != nil {
			return err
		}
		fmt.Printf("Created function: %s\n", lambdaConfig.Name)

	} else {
		fmt.Printf("Updating function: %s\n", lambdaConfig.Name)
		_, err := ggSession.lambda.UpdateFunctionCode(&lambda.UpdateFunctionCodeInput{
			FunctionName: &lambdaConfig.Name,
			ZipFile:      zip.Bytes(),
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated function: %s\n", lambdaConfig.Name)

	}

	publishVersion, err := ggSession.lambda.PublishVersion(&lambda.PublishVersionInput{
		FunctionName: &name,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Published %s version %s\n", name, *publishVersion.Version)

	// Create or update alias
	_, err = ggSession.lambda.GetAlias(&lambda.GetAliasInput{
		FunctionName: &name,
		Name:         &latest,
	})

	if err != nil {
		description := fmt.Sprintf("%s alias", lambdaConfig.Name)
		alias, err := ggSession.lambda.CreateAlias(&lambda.CreateAliasInput{
			Description:     &description,
			FunctionName:    publishVersion.FunctionName,
			FunctionVersion: publishVersion.Version,
			Name:            &latest,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Created alias: %v", alias)
		functionArn = *alias.AliasArn
	} else {
		alias, err := ggSession.lambda.UpdateAlias(&lambda.UpdateAliasInput{
			FunctionName:    publishVersion.FunctionName,
			FunctionVersion: publishVersion.Version,
			Name:            &latest,
		})
		if err != nil {
			return err
		}
		fmt.Printf("Updated alias: %v\n", alias)
		functionArn = *alias.AliasArn
	}
	ggSession.config.LambdaFunctions[name] = function{Arn: functionArn, ArnQualifier: latest}

	var functions []*greengrass.Function
	for k, v := range ggSession.config.LambdaFunctions {
		f, err := ggSession.lambda.GetFunction(&lambda.GetFunctionInput{
			FunctionName: &k,
			Qualifier:    &v.ArnQualifier,
		})
		if err != nil {
			continue
		}
		memorySize := *f.Configuration.MemorySize * 1000
		functions = append(functions, &greengrass.Function{
			FunctionArn: f.Configuration.FunctionArn,
			FunctionConfiguration: &greengrass.FunctionConfiguration{
				Executable: f.Configuration.Handler,
				MemorySize: &memorySize,
				Pinned:     &lambdaConfig.Pinned,
				Timeout:    f.Configuration.Timeout,
			},
			Id: &name,
		})
	}

	fmt.Printf("Starting CreateFunctionDefinition: %s\n", functionArn)

	if ggSession.config.FunctionDefinition.ID == "" {
		functionDefinition, err := ggSession.greengrass.CreateFunctionDefinition(&greengrass.CreateFunctionDefinitionInput{
			Name: &name,
		})
		if err != nil {
			return err
		}
		ggSession.config.FunctionDefinition.ID = *functionDefinition.Id
	}

	functionDefinitionVersion, err := ggSession.greengrass.CreateFunctionDefinitionVersion(&greengrass.CreateFunctionDefinitionVersionInput{
		FunctionDefinitionId: &ggSession.config.FunctionDefinition.ID,
		Functions:            functions,
	})
	if err != nil {
		return err
	}
	fmt.Printf("Created FunctionDefinitionVersion: %s\n", *functionDefinitionVersion.Arn)
	ggSession.config.FunctionDefinition.VersionArn = *functionDefinitionVersion.Arn

	ggSession.updateGroup()

	return nil
}

// ListLambda - list the lamba functions in the region
func (ggSession *GreengrassSession) ListLambda() error {
	functions, err := ggSession.greengrass.ListFunctionDefinitions(&greengrass.ListFunctionDefinitionsInput{})
	if err != nil {
		return err
	}
	for _, definition := range functions.Definitions {
		function, err := ggSession.greengrass.GetFunctionDefinitionVersion(&greengrass.GetFunctionDefinitionVersionInput{
			FunctionDefinitionId:        definition.Id,
			FunctionDefinitionVersionId: definition.LatestVersion,
		})
		if err != nil {
			return err
		}

		fmt.Printf("function %v\n", function)
	}

	return nil
}
