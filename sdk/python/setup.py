from setuptools import setup, Extension, find_packages

long_description = open('README.md', 'r').read()


setup(
        name='approzium',
        description='Python SDK for Approzium',
        long_description=long_description,
        long_description_content_type='text/markdown',
        url='https://github.com/approzium/approzium/',
        packages=find_packages(),
        ext_modules=[Extension('approzium._psycopg2_utils', ['approzium/_psycopg2_utils.c'])],
        classifiers=[
            "Programming Language :: Python :: 3",
            "Operating System :: OS Independent",
        ],
        install_requires=[
            'psycopg2',
            'boto3',
            'grpcio'
        ]
)
