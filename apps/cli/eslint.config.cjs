const js = require("@eslint/js");
const nx = require("@nx/eslint-plugin");
const typescriptParser = require("@typescript-eslint/parser");
const typescriptPlugin = require("@typescript-eslint/eslint-plugin");

module.exports = [
    js.configs.recommended,
    {
        files: ["**/*.ts", "**/*.js"],
        languageOptions: {
            parser: typescriptParser,
            parserOptions: {
                ecmaVersion: 2022,
                sourceType: "module",
            },
            globals: {
                console: "readonly",
                process: "readonly",
                Buffer: "readonly",
                __dirname: "readonly",
                __filename: "readonly",
                module: "readonly",
                require: "readonly",
                exports: "readonly",
                global: "readonly",
                setTimeout: "readonly",
                clearTimeout: "readonly",
                setInterval: "readonly",
                clearInterval: "readonly",
            },
        },
        plugins: {
            "@typescript-eslint": typescriptPlugin,
            "@nx": nx,
        },
        rules: {
            ...typescriptPlugin.configs.recommended.rules,
            "@nx/enforce-module-boundaries": ["error", {
                "enforceBuildableLibDependency": true,
                "allow": [],
                "depConstraints": [
                    {
                        "sourceTag": "*",
                        "onlyDependOnLibsWithTags": ["*"]
                    }
                ]
            }],
        },
    },
    {
        files: ["eslint.config.cjs"],
        languageOptions: {
            globals: {
                module: "readonly",
                require: "readonly",
                __dirname: "readonly",
            },
        },
        rules: {
            "no-unused-vars": "off",
        },
    },
];