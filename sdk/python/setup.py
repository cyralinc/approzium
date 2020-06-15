import setuptools

long_description = open('README.md', 'r').read()


setuptools.setup(
        name='approzium',
        description='Python SDK for Approzium',
        long_description=long_description,
        long_description_content_type='text/markdown',
        url='https://github.com/approzium/approzium/',
        packages=setuptools.find_packages(),
        classifiers=[
            "Programming Language :: Python :: 3",
            "Operating System :: OS Independent",
        ],
        install_requires=[
            'psycopg2',
            'boto3',
            'grpcio',
            'cython'
        ]
)
