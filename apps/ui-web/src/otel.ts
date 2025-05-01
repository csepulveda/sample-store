import { NodeSDK } from '@opentelemetry/sdk-node'
import { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
const { resourceFromAttributes } = require('@opentelemetry/resources');
const { SEMRESATTRS_SERVICE_NAME, SEMRESATTRS_DEPLOYMENT_ENVIRONMENT } = require('@opentelemetry/semantic-conventions');

const tempoUrl = process.env.OTEL_EXPORTER_OTLP_ENDPOINT || 'http://localhost:4318/v1/traces'

const sdk = new NodeSDK({
resource: resourceFromAttributes({
    [SEMRESATTRS_SERVICE_NAME]: 'ui-web',
    [SEMRESATTRS_DEPLOYMENT_ENVIRONMENT]: process.env.NODE_ENV || 'development',
    }),
  traceExporter: new OTLPTraceExporter({ url: tempoUrl }),
  instrumentations: [getNodeAutoInstrumentations()],
})

sdk.start()

process.on('SIGTERM', async () => {
  try {
    await sdk.shutdown()
    console.log('✅ OpenTelemetry SDK shut down gracefully')
  } catch (err: unknown) {
    console.error('❌ Error during shutdown', err)
  } finally {
    process.exit(0)
  }
})