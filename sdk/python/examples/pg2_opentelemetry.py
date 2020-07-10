from os import environ
import psycopg2
import approzium
import approzium.opentelemetry
from approzium.psycopg2 import connect
from opentelemetry import trace
from opentelemetry.ext import jaeger
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchExportSpanProcessor
from opentelemetry.ext.psycopg2 import Psycopg2Instrumentor

auth = approzium.AuthClient("authenticator:6000")
approzium.default_auth_client = auth

trace.set_tracer_provider(TracerProvider())

jaeger_exporter = jaeger.JaegerSpanExporter(
    service_name='approzium_service',
    agent_host_name='jaeger',
    agent_port=6831
)

trace.get_tracer_provider().add_span_processor(
    BatchExportSpanProcessor(jaeger_exporter)
)

tracer = trace.get_tracer(__name__)

# approzium.opentelemetry.instrument()
Psycopg2Instrumentor().instrument()

cnx = approzium.psycopg2.connect("host=dbmd5 dbname=db user=bob")
cursor = cnx.cursor()
with tracer.start_as_current_span("bar"):
    with tracer.start_as_current_span("baz"):
        print('Hello world!')
        cursor.execute("SELECT 1")
cursor.close()
cnx.close()
