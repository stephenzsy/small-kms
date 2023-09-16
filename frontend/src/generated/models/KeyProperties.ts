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
import type { CurveName } from './CurveName';
import {
    CurveNameFromJSON,
    CurveNameFromJSONTyped,
    CurveNameToJSON,
} from './CurveName';
import type { KeyType } from './KeyType';
import {
    KeyTypeFromJSON,
    KeyTypeFromJSONTyped,
    KeyTypeToJSON,
} from './KeyType';

/**
 * 
 * @export
 * @interface KeyProperties
 */
export interface KeyProperties {
    /**
     * 
     * @type {KeyType}
     * @memberof KeyProperties
     */
    kty: KeyType;
    /**
     * 
     * @type {number}
     * @memberof KeyProperties
     */
    keySize?: KeyPropertiesKeySizeEnum;
    /**
     * 
     * @type {CurveName}
     * @memberof KeyProperties
     */
    crv?: CurveName;
    /**
     * Keep using the same key version if exists
     * @type {boolean}
     * @memberof KeyProperties
     */
    reuseKey?: boolean;
}


/**
 * @export
 */
export const KeyPropertiesKeySizeEnum = {
    KeySize_2048: 2048,
    KeySize_3072: 3072,
    KeySize_4096: 4096
} as const;
export type KeyPropertiesKeySizeEnum = typeof KeyPropertiesKeySizeEnum[keyof typeof KeyPropertiesKeySizeEnum];


/**
 * Check if a given object implements the KeyProperties interface.
 */
export function instanceOfKeyProperties(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "kty" in value;

    return isInstance;
}

export function KeyPropertiesFromJSON(json: any): KeyProperties {
    return KeyPropertiesFromJSONTyped(json, false);
}

export function KeyPropertiesFromJSONTyped(json: any, ignoreDiscriminator: boolean): KeyProperties {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'kty': KeyTypeFromJSON(json['kty']),
        'keySize': !exists(json, 'key_size') ? undefined : json['key_size'],
        'crv': !exists(json, 'crv') ? undefined : CurveNameFromJSON(json['crv']),
        'reuseKey': !exists(json, 'reuse_key') ? undefined : json['reuse_key'],
    };
}

export function KeyPropertiesToJSON(value?: KeyProperties | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'kty': KeyTypeToJSON(value.kty),
        'key_size': value.keySize,
        'crv': CurveNameToJSON(value.crv),
        'reuse_key': value.reuseKey,
    };
}

