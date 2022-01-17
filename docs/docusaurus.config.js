module.exports = {
  title: 'Approzium',
  tagline: 'The tagline of my site',
  url: 'https://approzium.com',
  baseUrl: '/',
  favicon: 'img/apzm-icon.svg',
  organizationName: 'cyralinc',
  projectName: 'approzium',

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl: 'https://github.com/cyralinc/approzium/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],

  themeConfig: {
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
      algoliaOptions: {}, // Optional, if provided by Algolia
    },
    footer: {
      style: 'light',
      copyright: `Copyright © ${new Date().getFullYear()} Cyral, Inc.`,
    },
  },

  plugins: [
    '@docusaurus/plugin-ideal-image',
    [
      '@docusaurus/plugin-google-analytics',
      {
        trackingID: 'UA-172971690-1',
        anonymizeIP: true, // Should IPs be anonymized?
      },
    ],
  ],
}
