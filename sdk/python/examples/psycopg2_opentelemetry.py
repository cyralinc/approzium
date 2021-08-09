"""This is an example that shows Approzium Opentelemetry integration. It also
integrates with a Jaeger service to export and view generated traces.
"""
from opentelemetry import trace
from opentelemetry.exporter.jaeger.thrift import JaegerExporter
from opentelemetry.instrumentation.psycopg2 import Psycopg2Instrumentor
from opentelemetry.sdk.resources import SERVICE_NAME, Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.sdk.trace.export import BatchSpanProcessor

import approzium
import approzium.opentelemetry
from approzium.psycopg2 import connect

auth = approzium.AuthClient("authenticator:6001")
approzium.default_auth_client = auth

trace.set_tracer_provider(
    TracerProvider(resource=Resource.create({SERVICE_NAME: "approzium_service"}))
)
tracer = trace.get_tracer(__name__)

jaeger_exporter = JaegerExporter(agent_host_name="jaeger", agent_port=6831)

trace.get_tracer_provider().add_span_processor(BatchSpanProcessor(jaeger_exporter))

approzium.opentelemetry.instrument()
Psycopg2Instrumentor().instrument()

cnx = connect("host=dbmd5 dbname=db user=bob")
cursor = cnx.cursor()
with tracer.start_as_current_span("foo"):
    with tracer.start_as_current_span("bar"):
        print("Hello world!")
        cursor.execute("SELECT 1")
cursor.close()
cnx.close()
