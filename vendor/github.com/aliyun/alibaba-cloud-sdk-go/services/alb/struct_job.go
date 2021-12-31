package alb

//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.
//
// Code generated by Alibaba Cloud SDK Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

// Job is a nested struct in alb response
type Job struct {
	ApiName      string `json:"ApiName" xml:"ApiName"`
	CreateTime   int64  `json:"CreateTime" xml:"CreateTime"`
	ErrorCode    string `json:"ErrorCode" xml:"ErrorCode"`
	ErrorMessage string `json:"ErrorMessage" xml:"ErrorMessage"`
	Id           string `json:"Id" xml:"Id"`
	ModifyTime   int64  `json:"ModifyTime" xml:"ModifyTime"`
	OperateType  string `json:"OperateType" xml:"OperateType"`
	ResourceId   string `json:"ResourceId" xml:"ResourceId"`
	ResourceType string `json:"ResourceType" xml:"ResourceType"`
	Status       string `json:"Status" xml:"Status"`
}
