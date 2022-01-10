/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package gojmx

import (
	"context"
	"fmt"
	"github.com/newrelic/nrjmx/gojmx/internal/nrjmx"
	"github.com/newrelic/nrjmx/gojmx/internal/testutils"
	gopsutil "github.com/shirou/gopsutil/v3/process"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var timeStamp = time.Date(2022, time.January, 1, 01, 23, 45, 0, time.Local).UnixNano() / 1000000

func init() {
	_ = os.Setenv("NR_JMX_TOOL", filepath.Join(testutils.PrjDir, "bin", "nrjmx"))
	//_ = os.Setenv("NRIA_NRJMX_DEBUG", "true")
}

func Test_Query_Success_LargeAmountOfData(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	var data []map[string]interface{}

	name := strings.Repeat("tomas", 100)

	for i := 0; i < 1500; i++ {
		data = append(data, map[string]interface{}{
			"name":        fmt.Sprintf("%s-%d", name, i),
			"doubleValue": 1.2,
			"floatValue":  2.2,
			"numberValue": 3,
			"boolValue":   true,
			"dateValue":   timeStamp,
		})
	}

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeansBatch(ctx, container, data)
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname: jmxHost,
		Port:     int32(jmxPort.Int()),
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND query returns at least 5Mb of data.
	mBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.NoError(t, err)

	var result []*AttributeResponse

	for _, mBeanName := range mBeanNames {
		mBeanAttrNames, err := client.GetMBeanAttributeNames(mBeanName)
		assert.NoError(t, err)

		for _, mBeanAttrName := range mBeanAttrNames {
			jmxAttrs, err := client.GetMBeanAttributes(mBeanName, mBeanAttrName)
			assert.NoError(t, err)
			result = append(result, jmxAttrs...)
		}

	}
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(fmt.Sprintf("%v", result)), 5*1024*1024)
}

func Test_Query_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2222222,
		"numberValue": 3,
		"boolValue":   true,
		"dateValue":   timeStamp,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}

	actualMBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name", "DateValue",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*AttributeResponse{
		{
			Name: "test:type=Cat,name=tomas,attr=NumberValue",

			ResponseType: ResponseTypeInt,
			IntValue:     3,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=FloatValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  2.222222,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=DateValue",

			ResponseType: ResponseTypeString,
			StringValue:  "Jan 1, 2022 1:23:45 AM",
		},
		{
			Name: "test:type=Cat,name=tomas,attr=BoolValue",

			ResponseType: ResponseTypeBool,
			BoolValue:    true,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=DoubleValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=Name",

			ResponseType: ResponseTypeString,
			StringValue:  "tomas",
		},
	}

	var actual []*AttributeResponse
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttrs, err := client.GetMBeanAttributes("test:type=Cat,name=tomas", mBeanAttrName)
		assert.NoError(t, err)
		actual = append(actual, jmxAttrs...)
	}
	assert.ElementsMatch(t, expected, actual)
}

func Test_QueryMBean_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2222222,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	actualResponse, err := client.QueryMBeanAttributes("test:type=Cat,name=tomas")
	require.NoError(t, err)

	expectedResponse := []*AttributeResponse{
		{
			Name: "test:type=Cat,name=tomas,attr=FloatValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  2.222222,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=NumberValue",

			ResponseType: ResponseTypeInt,
			IntValue:     3,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=BoolValue",

			ResponseType: ResponseTypeBool,
			BoolValue:    true,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=DoubleValue",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name: "test:type=Cat,name=tomas,attr=Name",

			ResponseType: ResponseTypeString,
			StringValue:  "tomas",
		},
		{
			Name:         "test:type=Cat,name=tomas,attr=DateValue",
			ResponseType: ResponseTypeErr,
			StatusMsg:    `can't parse attribute, error: found a null value for bean: test:type=Cat,name=tomas,attr=DateValue, cause: null, stacktrace: null`,
		},
	}

	assert.ElementsMatch(t, expectedResponse, actualResponse)
}

func Test_Query_CompositeData(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMCompositeDataBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))

	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=CompositeDataCat,name=tomas",
	}
	actualMBeanNames, err := client.QueryMBeanNames("test:type=CompositeDataCat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"CatInfo",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=CompositeDataCat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	// AND Query returns expected data
	expected := []*AttributeResponse{
		{
			Name: "test:type=CompositeDataCat,name=tomas,attr=CatInfo.Double",

			ResponseType: ResponseTypeDouble,
			DoubleValue:  1.2,
		},
		{
			Name: "test:type=CompositeDataCat,name=tomas,attr=CatInfo.Name",

			ResponseType: ResponseTypeString,
			StringValue:  "tomas",
		},
	}

	actual, err := client.GetMBeanAttributes("test:type=CompositeDataCat,name=tomas", "CatInfo")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expected, actual)
}

func Test_Query_Timeout(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection fails
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: 1,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NotNil(t, client)
	assert.Error(t, err)
	defer assertCloseClientError(t, client)

	// AND Query returns expected error
	actual, err := client.GetMBeanAttributeNames("*:*")
	assert.Nil(t, actual)
	assert.Error(t, err)
}

func Test_ConnectionURL_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.2,
		"numberValue": 3,
		"boolValue":   true,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:    fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	actual, err := client.GetMBeanAttributes("test:type=Cat,name=tomas", "FloatValue")
	assert.NoError(t, err)

	expected := []*AttributeResponse{
		{
			Name:         "test:type=Cat,name=tomas,attr=FloatValue",
			ResponseType: ResponseTypeDouble,
			DoubleValue:  2.2,
		},
	}

	assert.Equal(t, expected, actual)
}

func Test_JavaNotInstalledError(t *testing.T) {
	// GIVEN a wrong Java Home
	os.Setenv("NRIA_JAVA_HOME", "/wrong/path")
	defer os.Unsetenv("NRIA_JAVA_HOME")

	ctx := context.Background()

	config := &JMXConfig{
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}
	// THEN connect fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "/wrong/path/bin/java")

	// AND Query fails with expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, errProcessNotRunning)
}

func Test_WrongMBeanFormatError(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be oppened
	config := &JMXConfig{
		ConnectionURL:    fmt.Sprintf("service:jmx:rmi:///jndi/rmi://%s:%s/jmxrmi", jmxHost, jmxPort.Port()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("wrong_format")
	assert.Nil(t, actual)
	assert.EqualError(t, err, `jmx error: cannot parse MBean glob pattern: 'wrong_format', valid: 'DOMAIN:BEAN', cause: Key properties cannot be empty, stacktrace: `)
}

func Test_Wrong_Connection(t *testing.T) {
	ctx := context.Background()

	// GIVEN a wrong hostname and port
	config := &JMXConfig{
		Hostname:         "localhost",
		Port:             1234,
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;")
	defer assertCloseClientError(t, client)

	// AND query returns expected error
	assert.Contains(t, err.Error(), "Connection refused to host: localhost;") // TODO: fix this, doesn't return the correct error

	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_SSLQuery_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Populate the JMX Server with mbeans
	resp, err := testutils.AddMBeans(ctx, container, map[string]interface{}{
		"name":        "tomas",
		"doubleValue": 1.2,
		"floatValue":  2.222222,
		"numberValue": 3,
		"boolValue":   true,
		"dateValue":   1641429296000,
	})
	assert.NoError(t, err)
	assert.Equal(t, "ok!\n", string(resp))
	defer testutils.CleanMBeans(ctx, container)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN SSL JMX connection can be opened
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           testutils.JmxUsername,
		Password:           testutils.JmxPassword,
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   testutils.KeystorePassword,
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMBeanNames := []string{
		"test:type=Cat,name=tomas",
	}
	actualMBeanNames, err := client.QueryMBeanNames("test:type=Cat,*")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanNames, actualMBeanNames)

	expectedMBeanAttrNames := []string{
		"BoolValue", "FloatValue", "NumberValue", "DoubleValue", "Name", "DateValue",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("test:name=tomas,type=Cat")
	require.NoError(t, err)
	require.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)
}

func Test_Wrong_Credentials(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           "wrong_username",
		Password:           "wrong_password",
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   testutils.KeystorePassword,
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}

	// THEN open fails with expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Authentication failed! Invalid username or password")
	defer assertCloseClientError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_Wrong_Certificate_password(t *testing.T) {
	ctx := context.Background()

	// GIVEN an SSL JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainerSSL(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// WHEN wrong jmx username and password is provided
	config := &JMXConfig{
		Hostname:           jmxHost,
		Port:               int32(jmxPort.Int()),
		Username:           testutils.JmxUsername,
		Password:           testutils.JmxPassword,
		KeyStore:           testutils.KeystorePath,
		KeyStorePassword:   "wrong_password",
		TrustStore:         testutils.TruststorePath,
		TrustStorePassword: testutils.TruststorePassword,
		RequestTimeoutMs:   testutils.DefaultTimeoutMs,
	}

	// THEN open returns expected error
	client, err := NewClient(ctx).Open(config)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SSLContext")
	defer assertCloseClientError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("test:type=Cat,*")
	assert.Nil(t, actual)
	assert.Errorf(t, err, "connection to JMX endpoint is not established")
}

func Test_Connector_Success(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JBoss Server with JMX exposed running inside a container
	container, err := testutils.RunJbossStandaloneJMXContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	// Install the connector
	dstFile := filepath.Join(testutils.PrjDir, "/bin/jboss-client.jar")
	err = testutils.CopyFileFromContainer(ctx, container, "/opt/jboss/wildfly/bin/client/jboss-client.jar", dstFile)
	assert.NoError(t, err)

	defer os.Remove(dstFile)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.JbossJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened
	config := &JMXConfig{
		Hostname:              jmxHost,
		Port:                  int32(jmxPort.Int()),
		Username:              testutils.JbossJMXUsername,
		Password:              testutils.JbossJMXPassword,
		IsJBossStandaloneMode: true,
		IsRemote:              true,
		RequestTimeoutMs:      testutils.DefaultTimeoutMs,
	}
	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)
	defer assertCloseClientNoError(t, client)

	// AND Query returns expected data
	expectedMbeanNames := []string{
		"jboss.as:subsystem=remoting,configuration=endpoint",
	}
	actualMbeanNames, err := client.QueryMBeanNames("jboss.as:subsystem=remoting,configuration=endpoint")
	assert.NoError(t, err)
	assert.ElementsMatch(t, expectedMbeanNames, actualMbeanNames)

	expectedMBeanAttrNames := []string{
		"authRealm",
		"authenticationRetries",
		"authorizeId",
		"bufferRegionSize",
		"heartbeatInterval",
		"maxInboundChannels",
		"maxInboundMessageSize",
		"maxInboundMessages",
		"maxOutboundChannels",
		"maxOutboundMessageSize",
		"maxOutboundMessages",
		"receiveBufferSize",
		"receiveWindowSize",
		"saslProtocol",
		"sendBufferSize",
		"serverName",
		"transmitWindowSize",
		"worker",
	}
	actualMBeanAttrNames, err := client.GetMBeanAttributeNames("jboss.as:subsystem=remoting,configuration=endpoint")
	assert.ElementsMatch(t, expectedMBeanAttrNames, actualMBeanAttrNames)

	expected := []*AttributeResponse{
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authRealm",
			StatusMsg:    "can't parse attribute, error: found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=authRealm, cause: null, stacktrace: null",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authenticationRetries",
			ResponseType: ResponseTypeInt,
			IntValue:     3,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=authorizeId",
			StatusMsg:    "can't parse attribute, error: found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=authorizeId, cause: null, stacktrace: null",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=bufferRegionSize",
			StatusMsg:    "can't parse attribute, error: found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=bufferRegionSize, cause: null, stacktrace: null",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=heartbeatInterval",
			ResponseType: ResponseTypeInt,
			IntValue:     60000,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundChannels",
			ResponseType: ResponseTypeInt,
			IntValue:     40,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessageSize",
			ResponseType: ResponseTypeInt,
			IntValue:     9223372036854775807,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxInboundMessages",
			ResponseType: ResponseTypeInt,
			IntValue:     80,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundChannels",
			ResponseType: ResponseTypeInt,
			IntValue:     40,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessageSize",
			ResponseType: ResponseTypeInt,
			IntValue:     9223372036854775807,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=maxOutboundMessages",
			ResponseType: ResponseTypeInt,
			IntValue:     65535,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveBufferSize",
			ResponseType: ResponseTypeInt,
			IntValue:     8192,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=serverName",
			StatusMsg:    "can't parse attribute, error: found a null value for bean: jboss.as:subsystem=remoting,configuration=endpoint,attr=serverName, cause: null, stacktrace: null",
			ResponseType: ResponseTypeErr,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=receiveWindowSize",
			ResponseType: ResponseTypeInt,
			IntValue:     131072,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=saslProtocol",
			ResponseType: ResponseTypeString,
			StringValue:  "remote",
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=sendBufferSize",
			ResponseType: ResponseTypeInt,
			IntValue:     8192,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=transmitWindowSize",
			ResponseType: ResponseTypeInt,
			IntValue:     131072,
		},
		{
			Name:         "jboss.as:subsystem=remoting,configuration=endpoint,attr=worker",
			ResponseType: ResponseTypeString,
			StringValue:  "default",
		},
	}

	var actual []*AttributeResponse
	for _, mBeanAttrName := range expectedMBeanAttrNames {
		jmxAttrs, err := client.GetMBeanAttributes("jboss.as:subsystem=remoting,configuration=endpoint", mBeanAttrName)
		if err != nil {
			continue
		}
		actual = append(actual, jmxAttrs...)
	}

	assert.ElementsMatch(t, expected, actual)
}

func TestClientClose(t *testing.T) {
	ctx := context.Background()

	// GIVEN a JMX Server running inside a container
	container, err := testutils.RunJMXServiceContainer(ctx)
	require.NoError(t, err)
	defer container.Terminate(ctx)

	jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
	require.NoError(t, err)

	// THEN JMX connection can be opened.
	config := &JMXConfig{
		Hostname:         jmxHost,
		Port:             int32(jmxPort.Int()),
		RequestTimeoutMs: testutils.DefaultTimeoutMs,
	}

	client, err := NewClient(ctx).Open(config)
	assert.NoError(t, err)

	// WHEN closing the client there is no error
	assertCloseClientNoError(t, client)

	// AND Query returns expected error
	actual, err := client.QueryMBeanNames("*:*")
	assert.Nil(t, actual)
	assert.ErrorIs(t, err, errProcessNotRunning)

	assert.True(t, client.nrJMXProcess.getOSProcessState().Success())
}

func TestProcessExits(t *testing.T) {
	// gojmx starts nrjmx bash script which stats a java process.
	// We want to make sure that if gojmx process dies, java process stops also.
	// To reproduce this scenario, we run the current test twice:
	// - subprocess SHOULD_RUN_EXIT = 1
	// - main test SHOULD_RUN_EXIT unset.
	if os.Getenv("IS_SUBPROCESS") == "1" {
		ctx := context.Background()

		// GIVEN a JMX Server running inside a container
		container, err := testutils.RunJMXServiceContainer(ctx)
		require.NoError(t, err)
		defer container.Terminate(ctx)

		jmxHost, jmxPort, err := testutils.GetContainerMappedPort(ctx, container, testutils.TestServerJMXPort)
		require.NoError(t, err)

		// THEN JMX connection can be opened
		config := &JMXConfig{
			Hostname:         jmxHost,
			Port:             int32(jmxPort.Int()),
			RequestTimeoutMs: testutils.DefaultTimeoutMs,
		}
		client, err := NewClient(ctx).Open(config)
		require.NoError(t, err)

		f, err := os.OpenFile(os.Getenv("TMP_FILE"), os.O_WRONLY|os.O_TRUNC, 0644)
		require.NoError(t, err)
		defer f.Close()
		_, err = fmt.Fprintf(f, "%d\n", client.nrJMXProcess.getPID())
		require.NoError(t, err)
		<-time.After(5 * time.Minute)
	}

	// Create a temporary file to communicate get the pid of the subprocess.
	tmpFile, err := ioutil.TempFile("", "nrjmx_pid_test")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	// Run again the current test function with IS_SUBPROCESS=1
	cmd := exec.Command(os.Args[0], "-test.run="+t.Name())
	cmd.Env = append(os.Environ(), "IS_SUBPROCESS=1", "TMP_FILE="+tmpFile.Name())

	stdErrBuffer := nrjmx.NewDefaultLimitedBuffer()
	stdOutBuffer := nrjmx.NewDefaultLimitedBuffer()
	cmd.Stderr = stdErrBuffer
	cmd.Stdout = stdOutBuffer
	err = cmd.Start()

	ctx, waitCancel := context.WithCancel(context.Background())
	go func() {
		err := cmd.Wait()
		if ctx.Err() == nil && err != nil {
			assert.NoError(t, err)
			panic(fmt.Errorf("stdout: %s\nstderr: %s", stdOutBuffer, stdErrBuffer))
		}
	}()

	// Get the pid from the subprocess.
	var pid int32
	require.Eventually(t, func() bool {
		var err error
		pid, err = testutils.ReadPidFile(tmpFile.Name())
		if err != nil {
			return false
		}
		return true
	}, 30*time.Second, 50*time.Millisecond)

	p, err := gopsutil.NewProcess(pid)
	assert.NoError(t, err)
	ch, err := p.Children()
	assert.NoError(t, err)

	// Stop listening for cmd.Wait() error. We want to kill the subprocess, at this point an error is expected.
	waitCancel()
	err = cmd.Process.Kill()
	require.NoError(t, err)

	// Check that java pid does not exist anymore.
	require.Eventually(t, func() bool {
		up, err := ch[0].IsRunning()
		if err != nil {
			return false
		}
		// assert is not running
		return !up
	}, 5*time.Second, 50*time.Millisecond)
}

func assertCloseClientNoError(t *testing.T, client *Client) {
	assert.NotNil(t, client)
	assert.NoError(t, client.Close())
}

func assertCloseClientError(t *testing.T, client *Client) {
	assert.NotNil(t, client)
	assert.Error(t, client.Close())
}
