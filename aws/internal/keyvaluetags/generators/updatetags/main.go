//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"go/format"
	"log"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

const filename = `update_tags_gen.go`

var serviceNames = []string{
	"accessanalyzer",
	"acm",
	"acmpca",
	"amplify",
	"apigateway",
	"apigatewayv2",
	"appmesh",
	"appstream",
	"appsync",
	"athena",
	"backup",
	"cloud9",
	"cloudfront",
	"cloudhsmv2",
	"cloudtrail",
	"cloudwatch",
	"cloudwatchevents",
	"cloudwatchlogs",
	"codecommit",
	"codedeploy",
	"codepipeline",
	"codestarnotifications",
	"cognitoidentity",
	"cognitoidentityprovider",
	"configservice",
	"databasemigrationservice",
	"dataexchange",
	"datapipeline",
	"datasync",
	"dax",
	"devicefarm",
	"directconnect",
	"directoryservice",
	"dlm",
	"docdb",
	"dynamodb",
	"ec2",
	"ecr",
	"ecs",
	"efs",
	"eks",
	"elasticache",
	"elasticbeanstalk",
	"elasticsearchservice",
	"elb",
	"elbv2",
	"emr",
	"firehose",
	"fsx",
	"gamelift",
	"glacier",
	"globalaccelerator",
	"glue",
	"guardduty",
	"greengrass",
	"imagebuilder",
	"iot",
	"iotanalytics",
	"iotevents",
	"kafka",
	"kinesis",
	"kinesisanalytics",
	"kinesisanalyticsv2",
	"kinesisvideo",
	"kms",
	"lambda",
	"licensemanager",
	"lightsail",
	"mediaconnect",
	"mediaconvert",
	"medialive",
	"mediapackage",
	"mediastore",
	"mq",
	"neptune",
	"networkmanager",
	"opsworks",
	"organizations",
	"pinpoint",
	"qldb",
	"quicksight",
	"ram",
	"rds",
	"redshift",
	"resourcegroups",
	"route53",
	"route53resolver",
	"sagemaker",
	"secretsmanager",
	"securityhub",
	"sfn",
	"sns",
	"sqs",
	"ssm",
	"storagegateway",
	"swf",
	"synthetics",
	"transfer",
	"waf",
	"wafregional",
	"wafv2",
	"workspaces",
}

type TemplateData struct {
	ServiceNames []string
}

func main() {
	// Always sort to reduce any potential generation churn
	sort.Strings(serviceNames)

	templateData := TemplateData{
		ServiceNames: serviceNames,
	}
	templateFuncMap := template.FuncMap{
		"ClientType":                      keyvaluetags.ServiceClientType,
		"TagFunction":                     keyvaluetags.ServiceTagFunction,
		"TagFunctionBatchSize":            keyvaluetags.ServiceTagFunctionBatchSize,
		"TagInputCustomValue":             keyvaluetags.ServiceTagInputCustomValue,
		"TagInputIdentifierField":         keyvaluetags.ServiceTagInputIdentifierField,
		"TagInputIdentifierRequiresSlice": keyvaluetags.ServiceTagInputIdentifierRequiresSlice,
		"TagInputResourceTypeField":       keyvaluetags.ServiceTagInputResourceTypeField,
		"TagInputTagsField":               keyvaluetags.ServiceTagInputTagsField,
		"TagPackage":                      keyvaluetags.ServiceTagPackage,
		"Title":                           strings.Title,
		"UntagFunction":                   keyvaluetags.ServiceUntagFunction,
		"UntagInputCustomValue":           keyvaluetags.ServiceUntagInputCustomValue,
		"UntagInputRequiresTagKeyType":    keyvaluetags.ServiceUntagInputRequiresTagKeyType,
		"UntagInputRequiresTagType":       keyvaluetags.ServiceUntagInputRequiresTagType,
		"UntagInputTagsField":             keyvaluetags.ServiceUntagInputTagsField,
	}

	tmpl, err := template.New("updatetags").Funcs(templateFuncMap).Parse(templateBody)

	if err != nil {
		log.Fatalf("error parsing template: %s", err)
	}

	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, templateData)

	if err != nil {
		log.Fatalf("error executing template: %s", err)
	}

	generatedFileContents, err := format.Source(buffer.Bytes())

	if err != nil {
		log.Fatalf("error formatting generated file: %s", err)
	}

	f, err := os.Create(filename)

	if err != nil {
		log.Fatalf("error creating file (%s): %s", filename, err)
	}

	defer f.Close()

	_, err = f.Write(generatedFileContents)

	if err != nil {
		log.Fatalf("error writing to file (%s): %s", filename, err)
	}
}

var templateBody = `
// Code generated by generators/updatetags/main.go; DO NOT EDIT.

package keyvaluetags

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
{{- range .ServiceNames }}
	"github.com/aws/aws-sdk-go/service/{{ . }}"
{{- end }}
)
{{ range .ServiceNames }}

// {{ . | Title }}UpdateTags updates {{ . }} service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func {{ . | Title }}UpdateTags(conn {{ . | ClientType }}, identifier string{{ if . | TagInputResourceTypeField }}, resourceType string{{ end }}, oldTagsMap interface{}, newTagsMap interface{}) error {
	oldTags := New(oldTagsMap)
	newTags := New(newTagsMap)
	{{- if eq (. | TagFunction) (. | UntagFunction) }}
	removedTags := oldTags.Removed(newTags)
	updatedTags := oldTags.Updated(newTags)

	// Ensure we do not send empty requests
	if len(removedTags) == 0 && len(updatedTags) == 0 {
		return nil
	}

	input := &{{ . | TagPackage }}.{{ . | TagFunction }}Input{
		{{- if . | TagInputIdentifierRequiresSlice }}
		{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
		{{- else }}
		{{ . | TagInputIdentifierField }}:   aws.String(identifier),
		{{- end }}
		{{- if . | TagInputResourceTypeField }}
		{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
		{{- end }}
	}

	if len(updatedTags) > 0 {
		input.{{ . | TagInputTagsField }} = updatedTags.IgnoreAws().{{ . | Title }}Tags()
	}

	if len(removedTags) > 0 {
		{{- if . | UntagInputRequiresTagType }}
		input.{{ . | UntagInputTagsField }} = removedTags.IgnoreAws().{{ . | Title }}Tags()
		{{- else if . | UntagInputRequiresTagKeyType }}
		input.{{ . | UntagInputTagsField }} = removedTags.IgnoreAws().{{ . | Title }}TagKeys()
		{{- else if . | UntagInputCustomValue }}
		input.{{ . | UntagInputTagsField }} = {{ . | UntagInputCustomValue }}
		{{- else }}
		input.{{ . | UntagInputTagsField }} = aws.StringSlice(removedTags.Keys())
		{{- end }}
	}

	_, err := conn.{{ . | TagFunction }}(input)

	if err != nil {
		return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
	}

	{{- else }}

	if removedTags := oldTags.Removed(newTags); len(removedTags) > 0 {
		{{- if . | TagFunctionBatchSize }}
		for _, removedTags := range removedTags.Chunks({{ . | TagFunctionBatchSize }}) {
		{{- end }}
		input := &{{ . | TagPackage }}.{{ . | UntagFunction }}Input{
			{{- if . | TagInputIdentifierRequiresSlice }}
			{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
			{{- else }}
			{{ . | TagInputIdentifierField }}:   aws.String(identifier),
			{{- end }}
			{{- if . | TagInputResourceTypeField }}
			{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
			{{- end }}
			{{- if . | UntagInputRequiresTagType }}
			{{ . | UntagInputTagsField }}:       removedTags.IgnoreAws().{{ . | Title }}Tags(),
			{{- else if . | UntagInputRequiresTagKeyType }}
			{{ . | UntagInputTagsField }}:       removedTags.IgnoreAws().{{ . | Title }}TagKeys(),
			{{- else if . | UntagInputCustomValue }}
			{{ . | UntagInputTagsField }}:       {{ . | UntagInputCustomValue }},
			{{- else }}
			{{ . | UntagInputTagsField }}:       aws.StringSlice(removedTags.IgnoreAws().Keys()),
			{{- end }}
		}

		_, err := conn.{{ . | UntagFunction }}(input)

		if err != nil {
			return fmt.Errorf("error untagging resource (%s): %w", identifier, err)
		}
		{{- if . | TagFunctionBatchSize }}
		}
		{{- end }}
	}

	if updatedTags := oldTags.Updated(newTags); len(updatedTags) > 0 {
		{{- if . | TagFunctionBatchSize }}
		for _, updatedTags := range updatedTags.Chunks({{ . | TagFunctionBatchSize }}) {
		{{- end }}
		input := &{{ . | TagPackage }}.{{ . | TagFunction }}Input{
			{{- if . | TagInputIdentifierRequiresSlice }}
			{{ . | TagInputIdentifierField }}:   aws.StringSlice([]string{identifier}),
			{{- else }}
			{{ . | TagInputIdentifierField }}:   aws.String(identifier),
			{{- end }}
			{{- if . | TagInputResourceTypeField }}
			{{ . | TagInputResourceTypeField }}: aws.String(resourceType),
			{{- end }}
			{{- if . | TagInputCustomValue }}
			{{ . | TagInputTagsField }}:         {{ . | TagInputCustomValue }},
			{{- else }}
			{{ . | TagInputTagsField }}:         updatedTags.IgnoreAws().{{ . | Title }}Tags(),
			{{- end }}
		}

		_, err := conn.{{ . | TagFunction }}(input)

		if err != nil {
			return fmt.Errorf("error tagging resource (%s): %w", identifier, err)
		}
		{{- if . | TagFunctionBatchSize }}
		}
		{{- end }}
	}

	{{- end }}

	return nil
}
{{- end }}
`
