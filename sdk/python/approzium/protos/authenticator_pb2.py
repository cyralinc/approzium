# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: authenticator.proto

from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import enum_type_wrapper

# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='authenticator.proto',
  package='approzium.authenticator.protos',
  syntax='proto3',
  serialized_options=None,
  serialized_pb=b'\n\x13\x61uthenticator.proto\x12\x1e\x61pprozium.authenticator.protos\"\x8f\x02\n\x10PGMD5HashRequest\x12:\n\x08\x61uthtype\x18\x01 \x01(\x0e\x32(.approzium.authenticator.protos.AuthType\x12G\n\x0f\x63lient_language\x18\x02 \x01(\x0e\x32..approzium.authenticator.protos.ClientLanguage\x12\x0e\n\x06\x64\x62host\x18\x03 \x01(\t\x12\x0e\n\x06\x64\x62port\x18\x04 \x01(\t\x12\x0e\n\x06\x64\x62user\x18\x05 \x01(\t\x12\x38\n\x07\x61wsauth\x18\x06 \x01(\x0b\x32\'.approzium.authenticator.protos.AWSAuth\x12\x0c\n\x04salt\x18\x07 \x01(\x0c\"\xc2\x02\n\x13PGSHA256HashRequest\x12:\n\x08\x61uthtype\x18\x01 \x01(\x0e\x32(.approzium.authenticator.protos.AuthType\x12G\n\x0f\x63lient_language\x18\x02 \x01(\x0e\x32..approzium.authenticator.protos.ClientLanguage\x12\x0e\n\x06\x64\x62host\x18\x03 \x01(\t\x12\x0e\n\x06\x64\x62port\x18\x04 \x01(\t\x12\x0e\n\x06\x64\x62user\x18\x05 \x01(\t\x12\x38\n\x07\x61wsauth\x18\x06 \x01(\x0b\x32\'.approzium.authenticator.protos.AWSAuth\x12\x0c\n\x04salt\x18\x07 \x01(\t\x12\x12\n\niterations\x18\x08 \x01(\r\x12\x1a\n\x12\x61uthentication_msg\x18\t \x01(\t\"\x1d\n\rPGMD5Response\x12\x0c\n\x04hash\x18\x01 \x01(\t\"2\n\x10PGSHA256Response\x12\x0e\n\x06\x63proof\x18\x01 \x01(\t\x12\x0e\n\x06sproof\x18\x02 \x01(\t\"F\n\x07\x41WSAuth\x12\"\n\x1asigned_get_caller_identity\x18\x01 \x01(\t\x12\x17\n\x0f\x63laimed_iam_arn\x18\x02 \x01(\t**\n\x08\x41uthType\x12\x15\n\x11TYPE_NOT_PROVIDED\x10\x00\x12\x07\n\x03\x41WS\x10\x01*?\n\x0e\x43lientLanguage\x12\x19\n\x15LANGUAGE_NOT_PROVIDED\x10\x00\x12\n\n\x06PYTHON\x10\x01\x12\x06\n\x02GO\x10\x02\x32\xfe\x01\n\rAuthenticator\x12q\n\x0cGetPGMD5Hash\x12\x30.approzium.authenticator.protos.PGMD5HashRequest\x1a-.approzium.authenticator.protos.PGMD5Response\"\x00\x12z\n\x0fGetPGSHA256Hash\x12\x33.approzium.authenticator.protos.PGSHA256HashRequest\x1a\x30.approzium.authenticator.protos.PGSHA256Response\"\x00\x62\x06proto3'
)

_AUTHTYPE = _descriptor.EnumDescriptor(
  name='AuthType',
  full_name='approzium.authenticator.protos.AuthType',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='TYPE_NOT_PROVIDED', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='AWS', index=1, number=1,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=809,
  serialized_end=851,
)
_sym_db.RegisterEnumDescriptor(_AUTHTYPE)

AuthType = enum_type_wrapper.EnumTypeWrapper(_AUTHTYPE)
_CLIENTLANGUAGE = _descriptor.EnumDescriptor(
  name='ClientLanguage',
  full_name='approzium.authenticator.protos.ClientLanguage',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='LANGUAGE_NOT_PROVIDED', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='PYTHON', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='GO', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=853,
  serialized_end=916,
)
_sym_db.RegisterEnumDescriptor(_CLIENTLANGUAGE)

ClientLanguage = enum_type_wrapper.EnumTypeWrapper(_CLIENTLANGUAGE)
TYPE_NOT_PROVIDED = 0
AWS = 1
LANGUAGE_NOT_PROVIDED = 0
PYTHON = 1
GO = 2



_PGMD5HASHREQUEST = _descriptor.Descriptor(
  name='PGMD5HashRequest',
  full_name='approzium.authenticator.protos.PGMD5HashRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='authtype', full_name='approzium.authenticator.protos.PGMD5HashRequest.authtype', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='client_language', full_name='approzium.authenticator.protos.PGMD5HashRequest.client_language', index=1,
      number=2, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbhost', full_name='approzium.authenticator.protos.PGMD5HashRequest.dbhost', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbport', full_name='approzium.authenticator.protos.PGMD5HashRequest.dbport', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbuser', full_name='approzium.authenticator.protos.PGMD5HashRequest.dbuser', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='awsauth', full_name='approzium.authenticator.protos.PGMD5HashRequest.awsauth', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='salt', full_name='approzium.authenticator.protos.PGMD5HashRequest.salt', index=6,
      number=7, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=b"",
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=56,
  serialized_end=327,
)


_PGSHA256HASHREQUEST = _descriptor.Descriptor(
  name='PGSHA256HashRequest',
  full_name='approzium.authenticator.protos.PGSHA256HashRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='authtype', full_name='approzium.authenticator.protos.PGSHA256HashRequest.authtype', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='client_language', full_name='approzium.authenticator.protos.PGSHA256HashRequest.client_language', index=1,
      number=2, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbhost', full_name='approzium.authenticator.protos.PGSHA256HashRequest.dbhost', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbport', full_name='approzium.authenticator.protos.PGSHA256HashRequest.dbport', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dbuser', full_name='approzium.authenticator.protos.PGSHA256HashRequest.dbuser', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='awsauth', full_name='approzium.authenticator.protos.PGSHA256HashRequest.awsauth', index=5,
      number=6, type=11, cpp_type=10, label=1,
      has_default_value=False, default_value=None,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='salt', full_name='approzium.authenticator.protos.PGSHA256HashRequest.salt', index=6,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='iterations', full_name='approzium.authenticator.protos.PGSHA256HashRequest.iterations', index=7,
      number=8, type=13, cpp_type=3, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='authentication_msg', full_name='approzium.authenticator.protos.PGSHA256HashRequest.authentication_msg', index=8,
      number=9, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=330,
  serialized_end=652,
)


_PGMD5RESPONSE = _descriptor.Descriptor(
  name='PGMD5Response',
  full_name='approzium.authenticator.protos.PGMD5Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='hash', full_name='approzium.authenticator.protos.PGMD5Response.hash', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=654,
  serialized_end=683,
)


_PGSHA256RESPONSE = _descriptor.Descriptor(
  name='PGSHA256Response',
  full_name='approzium.authenticator.protos.PGSHA256Response',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='cproof', full_name='approzium.authenticator.protos.PGSHA256Response.cproof', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='sproof', full_name='approzium.authenticator.protos.PGSHA256Response.sproof', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=685,
  serialized_end=735,
)


_AWSAUTH = _descriptor.Descriptor(
  name='AWSAuth',
  full_name='approzium.authenticator.protos.AWSAuth',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='signed_get_caller_identity', full_name='approzium.authenticator.protos.AWSAuth.signed_get_caller_identity', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='claimed_iam_arn', full_name='approzium.authenticator.protos.AWSAuth.claimed_iam_arn', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=b"".decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=737,
  serialized_end=807,
)

_PGMD5HASHREQUEST.fields_by_name['authtype'].enum_type = _AUTHTYPE
_PGMD5HASHREQUEST.fields_by_name['client_language'].enum_type = _CLIENTLANGUAGE
_PGMD5HASHREQUEST.fields_by_name['awsauth'].message_type = _AWSAUTH
_PGSHA256HASHREQUEST.fields_by_name['authtype'].enum_type = _AUTHTYPE
_PGSHA256HASHREQUEST.fields_by_name['client_language'].enum_type = _CLIENTLANGUAGE
_PGSHA256HASHREQUEST.fields_by_name['awsauth'].message_type = _AWSAUTH
DESCRIPTOR.message_types_by_name['PGMD5HashRequest'] = _PGMD5HASHREQUEST
DESCRIPTOR.message_types_by_name['PGSHA256HashRequest'] = _PGSHA256HASHREQUEST
DESCRIPTOR.message_types_by_name['PGMD5Response'] = _PGMD5RESPONSE
DESCRIPTOR.message_types_by_name['PGSHA256Response'] = _PGSHA256RESPONSE
DESCRIPTOR.message_types_by_name['AWSAuth'] = _AWSAUTH
DESCRIPTOR.enum_types_by_name['AuthType'] = _AUTHTYPE
DESCRIPTOR.enum_types_by_name['ClientLanguage'] = _CLIENTLANGUAGE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

PGMD5HashRequest = _reflection.GeneratedProtocolMessageType('PGMD5HashRequest', (_message.Message,), {
  'DESCRIPTOR' : _PGMD5HASHREQUEST,
  '__module__' : 'authenticator_pb2'
  # @@protoc_insertion_point(class_scope:approzium.authenticator.protos.PGMD5HashRequest)
  })
_sym_db.RegisterMessage(PGMD5HashRequest)

PGSHA256HashRequest = _reflection.GeneratedProtocolMessageType('PGSHA256HashRequest', (_message.Message,), {
  'DESCRIPTOR' : _PGSHA256HASHREQUEST,
  '__module__' : 'authenticator_pb2'
  # @@protoc_insertion_point(class_scope:approzium.authenticator.protos.PGSHA256HashRequest)
  })
_sym_db.RegisterMessage(PGSHA256HashRequest)

PGMD5Response = _reflection.GeneratedProtocolMessageType('PGMD5Response', (_message.Message,), {
  'DESCRIPTOR' : _PGMD5RESPONSE,
  '__module__' : 'authenticator_pb2'
  # @@protoc_insertion_point(class_scope:approzium.authenticator.protos.PGMD5Response)
  })
_sym_db.RegisterMessage(PGMD5Response)

PGSHA256Response = _reflection.GeneratedProtocolMessageType('PGSHA256Response', (_message.Message,), {
  'DESCRIPTOR' : _PGSHA256RESPONSE,
  '__module__' : 'authenticator_pb2'
  # @@protoc_insertion_point(class_scope:approzium.authenticator.protos.PGSHA256Response)
  })
_sym_db.RegisterMessage(PGSHA256Response)

AWSAuth = _reflection.GeneratedProtocolMessageType('AWSAuth', (_message.Message,), {
  'DESCRIPTOR' : _AWSAUTH,
  '__module__' : 'authenticator_pb2'
  # @@protoc_insertion_point(class_scope:approzium.authenticator.protos.AWSAuth)
  })
_sym_db.RegisterMessage(AWSAuth)



_AUTHENTICATOR = _descriptor.ServiceDescriptor(
  name='Authenticator',
  full_name='approzium.authenticator.protos.Authenticator',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=919,
  serialized_end=1173,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetPGMD5Hash',
    full_name='approzium.authenticator.protos.Authenticator.GetPGMD5Hash',
    index=0,
    containing_service=None,
    input_type=_PGMD5HASHREQUEST,
    output_type=_PGMD5RESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetPGSHA256Hash',
    full_name='approzium.authenticator.protos.Authenticator.GetPGSHA256Hash',
    index=1,
    containing_service=None,
    input_type=_PGSHA256HASHREQUEST,
    output_type=_PGSHA256RESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_AUTHENTICATOR)

DESCRIPTOR.services_by_name['Authenticator'] = _AUTHENTICATOR

# @@protoc_insertion_point(module_scope)
