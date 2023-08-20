module.exports = {
    root: true,
    extends: [
        'eslint:recommended',
        'plugin:@typescript-eslint/recommended',
        'plugin:svelte/recommended',
    ],
    parser: '@typescript-eslint/parser',
    plugins: ['@typescript-eslint'],
    parserOptions: {
        sourceType: 'module',
        ecmaVersion: 2020,
        extraFileExtensions: ['.svelte'],
        project: ['./tsconfig.json'],
    },
    ignorePatterns: [
        '.eslintrc.cjs',
        'src/app.html',
        'src/app.postcss',
        'src/lib/server/drizzle/**'
    ],
    rules: {
        semi: 'error',
        curly: 'error',
        'default-case': 'error',
        'default-param-last': 'error',
        eqeqeq: 'error',
        'max-depth': [
            'error',
            { max: 3 },
        ],
        'no-throw-literal': 'error',
        'no-return-assign': 'error',
        'no-sequences': 'error',
        'no-var': 'error',
        'prefer-const': 'error',
        'prefer-arrow-callback': 'error',
        'prefer-destructuring': 'error',
        'prefer-promise-reject-errors': 'error',
        'prefer-template': 'error',
        'require-await': 'error',
        yoda: 'error',
        'array-bracket-newline': [
            'error',
            { multiline: true, minItems: 2 },
        ],
        'array-bracket-spacing': [
            'error',
            'never',
            { arraysInArrays: false },
        ],
        'array-element-newline': [
            'error',
            { multiline: true, minItems: 2 },
        ],
        'object-curly-newline': [
            'error',
            { multiline: true, minProperties: 3 },
        ],
        'object-curly-spacing': [
            'error',
            'always',
        ],
        'arrow-parens': [
            'error',
            'always',
        ],
        'brace-style': [
            'error',
            '1tbs',
        ],
        'comma-dangle': [
            'error',
            'always-multiline',
        ],
        'eol-last': 'error',
        quotes: [
            'error',
            'single',
            {
                avoidEscape: true,
                allowTemplateLiterals: true,
            },
        ],
        indent: [
            'error',
            4,
        ],
        'func-call-spacing': [
            'error',
            'never',
        ],
        'max-len': [
            'warn',
            {
                code: 100,
                ignoreTrailingComments: true,
                ignoreUrls: true,
            },
        ],
        'newline-per-chained-call': [
            'error',
            { ignoreChainWithDepth: 2 },
        ],
        '@typescript-eslint/await-thenable': 'error',
        '@typescript-eslint/no-floating-promises': 'error',
        '@typescript-eslint/no-namespace': 'error',
        '@typescript-eslint/no-unsafe-argument': 'error',
        '@typescript-eslint/no-unsafe-assignment': 'error',
        'no-shadow': 'warn',
        'no-multi-spaces': 'warn',
        'no-multiple-empty-lines': 'warn',
        'no-trailing-spaces': 'warn',
        'max-lines-per-function': [
            'warn',
            { max: 20 },
        ],
        'no-console': 'warn',
        'default-case-last': 'warn',
        'class-methods-use-this': 'warn',
        '@typescript-eslint/no-explicit-any': 'warn',
        '@typescript-eslint/no-misused-promises': 'warn',
        '@typescript-eslint/no-non-null-asserted-optional-chain': 'warn',
    },
    env: {
        browser: true,
        es2017: true,
        node: true,
    },
    overrides: [
        {
            files: ['*.svelte'],
            parser: 'svelte-eslint-parser',
            parserOptions: { parser: '@typescript-eslint/parser' },
        },
    ],
};
