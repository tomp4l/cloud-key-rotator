// Copyright 2019 OVO Technology
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package location

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsSsm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/ovotech/cloud-key-rotator/pkg/cred"
)

// Ssm type
type Ssm struct {
	parameterName string
	region        string
}

func (ssm Ssm) Write(serviceAccountName string, keyWrapper KeyWrapper, creds cred.Credentials) (updated UpdatedLocation, err error) {
	var key string
	if key, err = getKeyForFileBasedLocation(keyWrapper); err != nil {
		return
	}
	svc := awsSsm.New(session.New())
	svc.Config.Region = aws.String(ssm.region)
	input := &awsSsm.PutParameterInput{
		Overwrite: aws.Bool(true),
		Name:      aws.String(ssm.parameterName),
		Value:     aws.String(key),
	}
	if _, err = svc.PutParameter(input); err != nil {
		return
	}
	updated = UpdatedLocation{
		LocationType: "SSM",
		LocationURI:  fmt.Sprintf("%s/%s", ssm.region, ssm.parameterName),
		LocationIDs:  []string{ssm.parameterName}}
	return
}
