// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github')
const darkCodeTheme = require('prism-react-renderer/themes/dracula')

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Approzium',
  tagline: 'The tagline of my site',
  url: 'https://approzium.com',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/apzm-icon.svg',
  organizationName: 'cyralinc',
  projectName: 'approzium',

  presets: [
    [
      'classic',
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/cyralinc/approzium/tree/main/docs',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
        googleAnalytics: {
          trackingID: 'UA-172971690-1',
          anonymizeIP: true, // Should IPs be anonymized?
        },
      }),
    ],
  ],
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      navbar: {
        title: 'Approzium',
        logo: {
          alt: 'Approzium Logo',
          src: 'img/apzm-icon.svg',
        },
        items: [
          {
            to: 'docs/',
            activeBasePath: 'docs',
            label: 'Documentation',
            position: 'right',
          },
          {
            href: 'https://github.com/cyralinc/approzium',
            label: 'GitHub',
            position: 'right',
          },
          {
            href: 'https://cyral.com',
            label: 'Cyral',
            position: 'right',
          },
        ],
      },
      algolia: {
        apiKey: '8ce076a38f6bf25b30877ca2251ec1e4',
        indexName: 'prod_APPROZIUM_DOCS',
        appId: '51LSIHDN1A', // Optional, if you run the DocSearch crawler on your own
      },
      footer: {
        style: 'light',
        copyright: `Copyright © ${new Date().getFullYear()} <a href="https://cyral.com/">Cyral, Inc.</a>`,
      },
      prism: {
        defaultLanguage: 'python',
      },
    }),

  plugins: ['docusaurus-plugin-sass', '@docusaurus/plugin-ideal-image'],
}

module.exports = config
