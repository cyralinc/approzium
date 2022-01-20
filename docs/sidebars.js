/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  someSidebar: [
    'overview',
    'quickstart',
    'architecture',
    'examples',
    'security-model',
    'observability',
    'installation',
    'configuration',
    'compatibility',
    {
      'Client Libraries': [
        {
          type: 'link',
          label: 'Python',
          href: 'https://approzium.readthedocs.io/en/latest/',
        },
      ],
    },
    'roadmap',
  ],
}

module.exports = sidebars
