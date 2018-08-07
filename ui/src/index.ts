// Polyfill at least promise - otherwise nothing works
import 'core-js/es6/promise';

async function commonEs6Features(): Promise<any> {
    if (
        'startsWith' in String.prototype &&
        'endsWith' in String.prototype &&
        'includes' in Array.prototype &&
        'assign' in Object &&
        'keys' in Object
    ) {
        return Promise.resolve();
    }

    return import('core-js');
}

async function main(): Promise<void> {
    await Promise.all([
        'fetch' in window ? Promise.resolve() : import('whatwg-fetch'),
        'Symbol' in window ? Promise.resolve() : import('core-js/es6/symbol'),
        'Map' in window ? Promise.resolve() : import('core-js/es6/map'),
        'Set' in window ? Promise.resolve() : import('core-js/es6/set'),
        commonEs6Features()
    ]);

    await import('src/main');
}

main();
