/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.0
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
/**
 * 
 * @export
 * @interface KeyParameters
 */
export interface KeyParameters {
    /**
     * 
     * @type {string}
     * @memberof KeyParameters
     */
    kty?: KeyParametersKtyEnum;
    /**
     * 
     * @type {number}
     * @memberof KeyParameters
     */
    size?: KeyParametersSizeEnum;
    /**
     * 
     * @type {string}
     * @memberof KeyParameters
     */
    curve?: KeyParametersCurveEnum;
}


/**
 * @export
 */
export const KeyParametersKtyEnum = {
    Kty_RSA: 'RSA',
    Kty_EC: 'EC'
} as const;
export type KeyParametersKtyEnum = typeof KeyParametersKtyEnum[keyof typeof KeyParametersKtyEnum];

/**
 * @export
 */
export const KeyParametersSizeEnum = {
    KeySize_2048: 2048,
    KeySize_3072: 3072,
    KeySize_4096: 4096
} as const;
export type KeyParametersSizeEnum = typeof KeyParametersSizeEnum[keyof typeof KeyParametersSizeEnum];

/**
 * @export
 */
export const KeyParametersCurveEnum = {
    EcCurve_P256: 'P-256',
    EcCurve_P384: 'P-384'
} as const;
export type KeyParametersCurveEnum = typeof KeyParametersCurveEnum[keyof typeof KeyParametersCurveEnum];


/**
 * Check if a given object implements the KeyParameters interface.
 */
export function instanceOfKeyParameters(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function KeyParametersFromJSON(json: any): KeyParameters {
    return KeyParametersFromJSONTyped(json, false);
}

export function KeyParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): KeyParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'kty': !exists(json, 'kty') ? undefined : json['kty'],
        'size': !exists(json, 'size') ? undefined : json['size'],
        'curve': !exists(json, 'curve') ? undefined : json['curve'],
    };
}

export function KeyParametersToJSON(value?: KeyParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'kty': value.kty,
        'size': value.size,
        'curve': value.curve,
    };
}

