module.exports = {
  title: 'Approzium',
  tagline: 'The tagline of my site',
  url: 'https://approzium.com',
  baseUrl: '/',
  favicon: 'img/apzm-icon.png',
  organizationName: 'cyralinc',
  projectName: 'approzium',
  themeConfig: {
    navbar: {
      title: 'Approzium',
      logo: {
        alt: 'Approzium Logo',
        src: 'img/apzm-icon.png',
      },
      links: [
        {
          to: 'docs/',
          activeBasePath: 'docs',
          label: 'Documentation',
          position: 'right',
        },
        {
          href: 'https://github.com/approzium/approzium',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'light',
      copyright: `Copyright © ${new Date().getFullYear()} Cyral, Inc.`,
    },
  },
  plugins: [
    '@docusaurus/plugin-ideal-image',
  ],
  presets: [
    [
      '@docusaurus/preset-classic',
      {
        docs: {
          // It is recommended to set document id as docs home page (`docs/` path).
          homePageId: 'overview',
          sidebarPath: require.resolve('./sidebars.js'),
          // Please change this to your repo.
          editUrl:
            'https://github.com/facebook/docusaurus/edit/master/website/',
        },
        theme: {
          customCss: require.resolve('./src/css/custom.css'),
        },
      },
    ],
  ],
};
