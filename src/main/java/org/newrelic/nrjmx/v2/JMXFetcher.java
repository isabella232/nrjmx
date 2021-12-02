package org.newrelic.nrjmx.v2;

/*
* Copyright 2020 New Relic Corporation. All rights reserved.
* SPDX-License-Identifier: Apache-2.0
*/

import java.io.IOException;
import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Properties;
import java.util.Set;
import java.util.logging.Logger;

import javax.management.Attribute;
import javax.management.InstanceNotFoundException;
import javax.management.IntrospectionException;
import javax.management.MBeanAttributeInfo;
import javax.management.MBeanInfo;
import javax.management.MBeanServerConnection;
import javax.management.MalformedObjectNameException;
import javax.management.ObjectInstance;
import javax.management.ObjectName;
import javax.management.ReflectionException;
import javax.management.openmbean.CompositeData;
import javax.management.remote.JMXConnector;
import javax.management.remote.JMXConnectorFactory;
import javax.management.remote.JMXServiceURL;
import javax.rmi.ssl.SslRMIClientSocketFactory;

import org.newrelic.nrjmx.v2.nrprotocol.JMXAttribute;
import org.newrelic.nrjmx.v2.nrprotocol.JMXConfig;
import org.newrelic.nrjmx.v2.nrprotocol.JMXConnectionError;
import org.newrelic.nrjmx.v2.nrprotocol.JMXError;
import org.newrelic.nrjmx.v2.nrprotocol.LogMessage;
import org.newrelic.nrjmx.v2.nrprotocol.ValueType;

/**
 * JMXFetcher class reads queries from an InputStream (usually stdin) and sends
 * the results to an OutputStream (usually stdout)
 */
public class JMXFetcher {
    public static final String defaultURIPath = "jmxrmi";
    public static final Boolean defaultJBossModeIsStandalone = false;

    private static final Logger logger = Logger.getLogger("nrjmx");

    private List<LogMessage> logs = new ArrayList<>();

    private MBeanServerConnection connection;
    private List<JMXAttribute> result = new ArrayList<>();
    private String connectionString;
    private Map<String, Object> connectionEnv = new HashMap<>();

    public class QueryError extends Exception {
        public QueryError(String message, Exception cause) {
            super(message, cause);
        }
    }

    public class ValueError extends Exception {
        public ValueError(String message) {
            super(message);
        }
    }

    public void connect() throws JMXConnectionError {
        try {
            JMXServiceURL address = new JMXServiceURL(connectionString);

            JMXConnector connector = JMXConnectorFactory.connect(address, connectionEnv);

            this.connection = connector.getMBeanServerConnection();
        } catch (Exception e) {
            String message = String.format("Can't connect to JMX server: '%s', error: '%s'", connectionString,
                    e.getMessage());
            throw new JMXConnectionError(1, message);
        }
    }

    public boolean StringIsNullOrEmpty(String value) {
        return value == null || value.equals("");
    }

    public JMXFetcher(JMXConfig jmxConfig) {
        if (jmxConfig.connectionURL != null && !jmxConfig.connectionURL.equals("")) {
            connectionString = jmxConfig.connectionURL;
          } else {
            // Official doc for remoting v3 is not available, see:
            // - https://developer.jboss.org/thread/196619
            // - http://jbossremoting.jboss.org/documentation/v3.html
            // Some doc on URIS at:
            // -
            // https://github.com/jboss-remoting/jboss-remoting/blob/master/src/main/java/org/jboss/remoting3/EndpointImpl.java#L292-L304
            // - https://stackoverflow.com/questions/42970921/what-is-http-remoting-protocol
            // -
            // http://www.mastertheboss.com/jboss-server/jboss-monitoring/using-jconsole-to-monitor-a-remote-wildfly-server
            String uriPath = jmxConfig.uriPath;
            if (jmxConfig.isRemote) {
              if (defaultURIPath.equals(uriPath)) {
                uriPath = "";
              } else {
                uriPath = uriPath.concat("/");
              }
              String remoteProtocol = "remote";
              if (jmxConfig.isJBossStandaloneMode) {
                remoteProtocol = "remote+http";
              }
              connectionString =
                  String.format("service:jmx:%s://%s:%s%s", remoteProtocol, jmxConfig.hostname, jmxConfig.port, uriPath);
            } else {
              connectionString =
                  String.format("service:jmx:rmi:///jndi/rmi://%s:%s/%s", jmxConfig.hostname, jmxConfig.port, uriPath);
            }
          }
      
          if (!"".equals(jmxConfig.username)) {
            connectionEnv.put(JMXConnector.CREDENTIALS, new String[] {jmxConfig.username, jmxConfig.password});
          }
      
          if (!"".equals(jmxConfig.keyStore) && !"".equals(jmxConfig.trustStore)) {
            Properties p = System.getProperties();
            p.put("javax.net.ssl.keyStore", jmxConfig.keyStore);
            p.put("javax.net.ssl.keyStorePassword", jmxConfig.keyStorePassword);
            p.put("javax.net.ssl.trustStore", jmxConfig.trustStore);
            p.put("javax.net.ssl.trustStorePassword", jmxConfig.trustStorePassword);
            connectionEnv.put("com.sun.jndi.rmi.factory.socket", new SslRMIClientSocketFactory());
          }
        }

    public List<LogMessage> getLogs() {
        List<LogMessage> result = this.logs;
        this.logs = new ArrayList<LogMessage>();
        return result;
    }

    public List<JMXAttribute> queryMbean(String beanName) throws JMXError {

        try {
            Set<ObjectInstance> beanInstances;
            // try {



            beanInstances = query(beanName);

            // } catch (JMXFetcher.QueryError e) {
            // // logger.warning(e.getMessage());
            // // logger.log(Level.FINE, e.getMessage(), e);
            // return null;
            // }

            for (ObjectInstance instance : beanInstances) {
                // try {
                queryAttributes(instance);

                // } catch (JMXFetcher.QueryError e) {
                // // logger.warning(e.getMessage());
                // // logger.log(Level.FINE, e.getMessage(), e);
                // }
            }
            try {

            return popResults();
            } catch (IllegalArgumentException e) {

            }
            return result;
        } catch (Exception ex) {
            throw new JMXError(ex.getMessage());
        }
        
    }

    private Set<ObjectInstance> query(String beanName) throws QueryError {
        ObjectName queryObject;

        try {
            queryObject = new ObjectName(beanName);
        } catch (MalformedObjectNameException e) {
            throw new QueryError("Can't parse bean name " + beanName, e);
        }

        Set<ObjectInstance> beanInstances;
        try {
            beanInstances = connection.queryMBeans(queryObject, null);
        } catch (IOException e) {
            throw new QueryError("Can't get beans for query " + beanName, e);
        }

        return beanInstances;
    }

    private void queryAttributes(ObjectInstance instance) throws QueryError, JMXError {
        ObjectName objectName = instance.getObjectName();
        MBeanInfo info;

        try {
            info = connection.getMBeanInfo(objectName);
        } catch (InstanceNotFoundException | IntrospectionException | ReflectionException | IOException e) {
            throw new QueryError("Can't find bean " + objectName.toString(), e);
        }

        MBeanAttributeInfo[] attrInfo = info.getAttributes();

        for (MBeanAttributeInfo attr : attrInfo) {
            if (!attr.isReadable()) {
                continue;
            }

            String attrName = attr.getName();
            Object value;

        
            try {
                value = connection.getAttribute(objectName, attrName);
                if (value instanceof Attribute) {
                    Attribute jmxAttr = (Attribute) value;
                    value = jmxAttr.getValue();
                }
            } catch (Exception e) {
                // logger.warning("Can't get attribute " + attrName + " for bean " +
                // objectName.toString() + ": "
                // + e.getMessage());
                continue;
            }

            String name = String.format("%s,attr=%s", objectName.toString(), attrName);
            try {
                parseValue(name, value);
            } catch (ValueError e) {
                this.logs.add(new LogMessage(e.getMessage()));
                // logger.fine(e.getMessage());
            }
        }
    }

    private List<JMXAttribute> popResults() {
        List<JMXAttribute> out = result;
        result = new ArrayList<>();
        return out;
    }

    private void parseValue(String name, Object value) throws ValueError {
        JMXAttribute attr = new JMXAttribute();
        attr.attribute = name;

        if (value == null) {
            throw new ValueError("Found a null value for bean " + name);
        } else if (value instanceof java.lang.Double) {
            attr.doubleValue = parseDouble((Double) value);
            attr.valueType = ValueType.DOUBLE;
            result.add(attr);
        } else if (value instanceof java.lang.Float) {
            attr.doubleValue = parseFloatToDouble((Float) value);
            attr.valueType = ValueType.DOUBLE;
            result.add(attr);
        } else if (value instanceof Number) {
            attr.intValue = ((Number) value).longValue();
            attr.valueType = ValueType.INT;
            result.add(attr);
        } else if (value instanceof String) {
            attr.stringValue = (String) value;
            attr.valueType = ValueType.STRING;
            result.add(attr);
        } else if (value instanceof Boolean) {
            attr.boolValue = (Boolean) value;
            attr.valueType = ValueType.BOOL;
            result.add(attr);
        } else if (value instanceof CompositeData) {
            CompositeData cdata = (CompositeData) value;
            Set<String> fieldKeys = cdata.getCompositeType().keySet();

            for (String field : fieldKeys) {
                if (field.length() < 1)
                    continue;

                String fieldKey = field.substring(0, 1).toUpperCase() + field.substring(1);
                parseValue(String.format("%s.%s", name, fieldKey), cdata.get(field));
            }
        } else if (value instanceof HashMap) {
            // TODO: Process hashmaps
            // logger.fine("HashMaps are not supported yet: " + name);
        } else if (value instanceof ArrayList || value.getClass().isArray()) {
            // TODO: Process arrays
            // logger.fine("Arrays are not supported yet: " + name);
        } else {
            throw new ValueError("Unsuported data type (" + value.getClass() + ") for bean " + name);
        }
    }

    /**
     * XXX: JSON does not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
     */
    private Double parseDouble(Double value) {
        if (value.isNaN()) {
            return 0.0;
        } else if (value == Double.NEGATIVE_INFINITY) {
            return Double.MIN_VALUE;
        } else if (value == Double.POSITIVE_INFINITY) {
            return Double.MAX_VALUE;
        }

        return value;
    }

    /**
     * XXX: JSON does not support NaN, Infinity, or -Infinity as they come back from
     * JMX. So we parse them out to 0, Max Double, and Min Double respectively.
     */
    private Double parseFloatToDouble(Float value) {
        if (value.isNaN()) {
            return 0.0d;
        } else if (value == Double.NEGATIVE_INFINITY) {
            return Double.MIN_VALUE;
        } else if (value == Double.POSITIVE_INFINITY) {
            return Double.MAX_VALUE;
        }

        return new BigDecimal(value.toString()).doubleValue();
    }
}