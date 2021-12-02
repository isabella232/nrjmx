package org.newrelic.nrjmx.v2;

import java.util.List;

import org.apache.thrift.TException;
import org.apache.thrift.server.TServer;
import org.newrelic.nrjmx.v2.nrprotocol.*;

public class JMXServiceHandler implements JMXService.Iface {

    private JMXFetcher jmxFetcher;
    private TServer server;

    @Override
    public void ping() throws JMXError, TException {
        return;
    }

    @Override
    public void connect(JMXConfig config) throws TException {
        this.jmxFetcher = new JMXFetcher(config);
        jmxFetcher.connect();
    }

    @Override
    public List<JMXAttribute> queryMbean(String beanName) throws TException {
        return jmxFetcher.queryMbean(beanName);
    }

    @Override
    public void disconnect() throws TException {
        if (this.server == null) {
            throw new TException("cannot disconnect, server handler null");
        }
        this.server.stop();
    }

    @Override
    public List<LogMessage> getLogs() throws TException {
        return jmxFetcher.getLogs();
    }

    public void addServer(TServer server) {
        this.server = server;
    }
}