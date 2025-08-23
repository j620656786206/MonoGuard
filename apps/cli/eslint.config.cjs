const { FlatCompat } = require("@eslint/eslintrc");
const js = require("@eslint/js");
const { fixupConfigRules } = require("@eslint/compat");
const nx = require("@nx/eslint-plugin");

const compat = new FlatCompat({
  baseDirectory: __dirname,
  recommendedConfig: js.configs.recommended,
});

module.exports = [
    ...fixupConfigRules(compat.extends("../../.eslintrc.json")),
    ...nx.configs["flat/typescript"],
    {
        files: [
            "**/*.ts",
            "**/*.js"
        ],
        env: {
            node: true,
        },
        // Override or add rules here
        rules: {}
    }
];