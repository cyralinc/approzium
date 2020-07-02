import opentelemetry.ext.dbapi
from opentelemetry.ext.dbapi import DatabaseApiIntegration


original_DatabaseApiIntegration = None

class ApproziumDatabaseApiIntegration(DatabaseApiIntegration):
    def get_connection_attributes(self, connection):
        super().get_connection_attributes(connection)
        # check if this is an Approzium connection
        if hasattr(connection, 'authenticator'):
            self.add_approzium_span_attributes(connection)

    def add_approzium_span_attributes(self, connection):
        authenticator = connection.authenticator
        for key, value in connection.authenticator.attribution_info.items():
            self.span_attributes['approzium.' + key] = value


def instrument():
    global original_DatabaseApiIntegration
    original_DatabaseApiIntegration = DatabaseApiIntegration
    opentelemetry.ext.dbapi.DatabaseApiIntegration = ApproziumDatabaseApiIntegration


def uninstrument():
    if original_DatabaseApiIntegration is None:
        raise Exception('Opentelemetry is not currently instrumented')
    opentelemetry.ext.dbapi.DatabaseApiIntegration = \
    original_DatabaseApiIntegration
