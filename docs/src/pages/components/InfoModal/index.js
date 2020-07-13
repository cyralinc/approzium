import React, { useState, useEffect } from 'react'
import clsx from 'clsx'
import CodeBlock from '@theme/CodeBlock'
import styles from './styles.module.css'


const info = {
  'observability': {
    description: 'Approzium provides rich execution context derived from cloud infrastructure metadata services, thus eliminating blind spots in the diagnosis and tracing of complex performance problems within autoscaled microservice environments.',
    codeHighlight: '{15-22}',
  },
  'security': {
    description: 'Approzium enables applications to connect to databases without requiring access to credentials, thus preventing leaks through inadvertent application logging, application compromise, or theft of secrets manager API keys.',
    codeHighlight: '{7-9}',
  }
}

const codeSample = `from approzium import AuthClient
from approzium.psycopg2 import connect

# create an Authenticator client
auth = AuthClient('authenticator:6001')

# connect using Approzium's connect method without a password
conn = connect("host=1.2.3.4 user=dbuser1 dbname=mydbhost",
               authenticator=auth)

# use the connection as you typically would. very cool!
cur = conn.cursor()
cur.execute('SELECT 1')

# attribute your database connections to a unique identity
auth.attribution_info()
# {
#     'authenticator_address': 'localhost:6001',
#     'iam_arn': '<verified IAM Amazon Resource Number>',
#     'authenticated': true,
#     'num_connections': 1,
# }
`

function InfoModal() {
  const [selectedCategory, setSelectedCategory] = useState('observability')

  const selectObservability = () => setSelectedCategory('observability')
  const selectSecurity = () => setSelectedCategory('security')

  const selectedInfo = info[selectedCategory]

  return (
    <div className={styles.container}>
      <div className={styles.textModal}>
        <div className={styles.selectorContainer}>
          <button
            className={clsx(styles.selector, styles.leftSelector, {
              [styles.selected]: selectedCategory == 'observability',
            })}
            onClick={selectObservability}>
            Observability
          </button>
          <button
            className={clsx(styles.selector, {
              [styles.selected]: selectedCategory == 'security',
            })}
            onClick={selectSecurity}>
            Security
          </button>
        </div>
        <p className={styles.description}>{selectedInfo.description}</p>
      </div>
      <CodeBlock className={'python'} metastring={selectedInfo.codeHighlight} landingPage={true}>
        {codeSample}
      </CodeBlock>
    </div>
  )
}

export default InfoModal;