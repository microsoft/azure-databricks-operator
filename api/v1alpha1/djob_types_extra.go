/*
The MIT License (MIT)

Copyright (c) 2019  Microsoft

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package v1alpha1

import (
	dbjobsmodels "github.com/polar-rams/databricks-sdk-golang/azure/jobs/models"
	dblibsmodels "github.com/polar-rams/databricks-sdk-golang/azure/libraries/models"
)

// JobSettings is similar to dbmodels.JobSettings, the reason it
// exists is because dbmodels.JobSettings doesn't support ExistingClusterName
// ExistingClusterName allows discovering databricks clusters by it's kubernetese object name
type JobSettings struct {
	ExistingClusterID      string                              `json:"existing_cluster_id,omitempty" url:"existing_cluster_id,omitempty"`
	ExistingClusterName    string                              `json:"existing_cluster_name,omitempty" url:"existing_cluster_name,omitempty"`
	NewCluster             *dbjobsmodels.NewCluster            `json:"new_cluster,omitempty" url:"new_cluster,omitempty"`
	NotebookTask           *dbjobsmodels.NotebookTask          `json:"notebook_task,omitempty" url:"notebook_task,omitempty"`
	SparkJarTask           *dbjobsmodels.SparkJarTask          `json:"spark_jar_task,omitempty" url:"spark_jar_task,omitempty"`
	SparkPythonTask        *dbjobsmodels.SparkPythonTask       `json:"spark_python_task,omitempty" url:"spark_python_task,omitempty"`
	SparkSubmitTask        *dbjobsmodels.SparkSubmitTask       `json:"spark_submit_task,omitempty" url:"spark_submit_task,omitempty"`
	Name                   string                              `json:"name,omitempty" url:"name,omitempty"`
	Libraries              []dblibsmodels.Library              `json:"libraries,omitempty" url:"libraries,omitempty"`
	EmailNotifications     *dbjobsmodels.JobEmailNotifications `json:"email_notifications,omitempty" url:"email_notifications,omitempty"`
	TimeoutSeconds         int32                               `json:"timeout_seconds,omitempty" url:"timeout_seconds,omitempty"`
	MaxRetries             int32                               `json:"max_retries,omitempty" url:"max_retries,omitempty"`
	MinRetryIntervalMillis int32                               `json:"min_retry_interval_millis,omitempty" url:"min_retry_interval_millis,omitempty"`
	RetryOnTimeout         bool                                `json:"retry_on_timeout,omitempty" url:"retry_on_timeout,omitempty"`
	Schedule               *dbjobsmodels.CronSchedule          `json:"schedule,omitempty" url:"schedule,omitempty"`
	MaxConcurrentRuns      int32                               `json:"max_concurrent_runs,omitempty" url:"max_concurrent_runs,omitempty"`
}

// ToK8sJobSettings converts a databricks JobSettings object to k8s JobSettings object.
// It is needed to add ExistingClusterName and follow k8s camleCase naming convention
func ToK8sJobSettings(dbjs *dbjobsmodels.JobSettings) JobSettings {
	var k8sjs JobSettings
	k8sjs.ExistingClusterID = dbjs.ExistingClusterID
	k8sjs.NewCluster = dbjs.NewCluster
	k8sjs.NotebookTask = dbjs.NotebookTask
	k8sjs.SparkJarTask = dbjs.SparkJarTask
	k8sjs.SparkPythonTask = dbjs.SparkPythonTask
	k8sjs.SparkSubmitTask = dbjs.SparkSubmitTask
	k8sjs.Name = dbjs.Name
	k8sjs.Libraries = *dbjs.Libraries
	k8sjs.EmailNotifications = dbjs.EmailNotifications
	k8sjs.TimeoutSeconds = dbjs.TimeoutSeconds
	k8sjs.MaxRetries = dbjs.MaxRetries
	k8sjs.MinRetryIntervalMillis = dbjs.MinRetryIntervalMillis
	k8sjs.RetryOnTimeout = dbjs.RetryOnTimeout
	k8sjs.Schedule = dbjs.Schedule
	k8sjs.MaxConcurrentRuns = dbjs.MaxConcurrentRuns
	return k8sjs
}

// ToDatabricksJobSettings converts a k8s JobSettings object to a DataBricks JobSettings object.
// It is needed to add ExistingClusterName and follow k8s camleCase naming convention
func ToDatabricksJobSettings(k8sjs *JobSettings) dbjobsmodels.JobSettings {

	var dbjs dbjobsmodels.JobSettings
	dbjs.ExistingClusterID = k8sjs.ExistingClusterID
	dbjs.NewCluster = k8sjs.NewCluster
	dbjs.NotebookTask = k8sjs.NotebookTask
	dbjs.SparkJarTask = k8sjs.SparkJarTask
	dbjs.SparkPythonTask = k8sjs.SparkPythonTask
	dbjs.SparkSubmitTask = k8sjs.SparkSubmitTask
	dbjs.Name = k8sjs.Name
	dbjs.Libraries = &k8sjs.Libraries
	dbjs.EmailNotifications = k8sjs.EmailNotifications
	dbjs.TimeoutSeconds = k8sjs.TimeoutSeconds
	dbjs.MaxRetries = k8sjs.MaxRetries
	dbjs.MinRetryIntervalMillis = k8sjs.MinRetryIntervalMillis
	dbjs.RetryOnTimeout = k8sjs.RetryOnTimeout
	dbjs.Schedule = k8sjs.Schedule
	dbjs.MaxConcurrentRuns = k8sjs.MaxConcurrentRuns
	return dbjs
}
