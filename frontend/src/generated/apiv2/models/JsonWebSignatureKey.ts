/* tslint:disable */
/* eslint-disable */
/**
 * Cryptocat API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.3
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import { exists, mapValues } from '../runtime';
import type { JsonWebKeyCurveName } from './JsonWebKeyCurveName';
import {
    JsonWebKeyCurveNameFromJSON,
    JsonWebKeyCurveNameFromJSONTyped,
    JsonWebKeyCurveNameToJSON,
} from './JsonWebKeyCurveName';
import type { JsonWebKeyOperation } from './JsonWebKeyOperation';
import {
    JsonWebKeyOperationFromJSON,
    JsonWebKeyOperationFromJSONTyped,
    JsonWebKeyOperationToJSON,
} from './JsonWebKeyOperation';
import type { JsonWebKeyType } from './JsonWebKeyType';
import {
    JsonWebKeyTypeFromJSON,
    JsonWebKeyTypeFromJSONTyped,
    JsonWebKeyTypeToJSON,
} from './JsonWebKeyType';

/**
 * 
 * @export
 * @interface JsonWebSignatureKey
 */
export interface JsonWebSignatureKey {
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    alg?: string;
    /**
     * 
     * @type {JsonWebKeyType}
     * @memberof JsonWebSignatureKey
     */
    kty: JsonWebKeyType;
    /**
     * 
     * @type {JsonWebKeyCurveName}
     * @memberof JsonWebSignatureKey
     */
    crv?: JsonWebKeyCurveName;
    /**
     * 
     * @type {number}
     * @memberof JsonWebSignatureKey
     */
    keySize?: number;
    /**
     * 
     * @type {Array<JsonWebKeyOperation>}
     * @memberof JsonWebSignatureKey
     */
    keyOps?: Array<JsonWebKeyOperation>;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    kid?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    n?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    e?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    x?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    y?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    x5t?: string;
    /**
     * 
     * @type {string}
     * @memberof JsonWebSignatureKey
     */
    x5tS256?: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof JsonWebSignatureKey
     */
    x5c?: Array<string>;
}

/**
 * Check if a given object implements the JsonWebSignatureKey interface.
 */
export function instanceOfJsonWebSignatureKey(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "kty" in value;

    return isInstance;
}

export function JsonWebSignatureKeyFromJSON(json: any): JsonWebSignatureKey {
    return JsonWebSignatureKeyFromJSONTyped(json, false);
}

export function JsonWebSignatureKeyFromJSONTyped(json: any, ignoreDiscriminator: boolean): JsonWebSignatureKey {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'alg': !exists(json, 'alg') ? undefined : json['alg'],
        'kty': JsonWebKeyTypeFromJSON(json['kty']),
        'crv': !exists(json, 'crv') ? undefined : JsonWebKeyCurveNameFromJSON(json['crv']),
        'keySize': !exists(json, 'key_size') ? undefined : json['key_size'],
        'keyOps': !exists(json, 'key_ops') ? undefined : ((json['key_ops'] as Array<any>).map(JsonWebKeyOperationFromJSON)),
        'kid': !exists(json, 'kid') ? undefined : json['kid'],
        'n': !exists(json, 'n') ? undefined : json['n'],
        'e': !exists(json, 'e') ? undefined : json['e'],
        'x': !exists(json, 'x') ? undefined : json['x'],
        'y': !exists(json, 'y') ? undefined : json['y'],
        'x5t': !exists(json, 'x5t') ? undefined : json['x5t'],
        'x5tS256': !exists(json, 'x5t#S256') ? undefined : json['x5t#S256'],
        'x5c': !exists(json, 'x5c') ? undefined : json['x5c'],
    };
}

export function JsonWebSignatureKeyToJSON(value?: JsonWebSignatureKey | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'alg': value.alg,
        'kty': JsonWebKeyTypeToJSON(value.kty),
        'crv': JsonWebKeyCurveNameToJSON(value.crv),
        'key_size': value.keySize,
        'key_ops': value.keyOps === undefined ? undefined : ((value.keyOps as Array<any>).map(JsonWebKeyOperationToJSON)),
        'kid': value.kid,
        'n': value.n,
        'e': value.e,
        'x': value.x,
        'y': value.y,
        'x5t': value.x5t,
        'x5t#S256': value.x5tS256,
        'x5c': value.x5c,
    };
}
