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
 * these attributes should mostly confirm to JWK (RFC7517)
 * @export
 * @interface JsonWebKeySpec
 */
export interface JsonWebKeySpec {
    /**
     * 
     * @type {string}
     * @memberof JsonWebKeySpec
     */
    alg?: string;
    /**
     * 
     * @type {JsonWebKeyType}
     * @memberof JsonWebKeySpec
     */
    kty?: JsonWebKeyType;
    /**
     * 
     * @type {JsonWebKeyCurveName}
     * @memberof JsonWebKeySpec
     */
    crv?: JsonWebKeyCurveName;
    /**
     * 
     * @type {number}
     * @memberof JsonWebKeySpec
     */
    keySize?: number;
    /**
     * 
     * @type {Array<JsonWebKeyOperation>}
     * @memberof JsonWebKeySpec
     */
    keyOps?: Array<JsonWebKeyOperation>;
    /**
     * 
     * @type {boolean}
     * @memberof JsonWebKeySpec
     */
    ext?: boolean;
}

/**
 * Check if a given object implements the JsonWebKeySpec interface.
 */
export function instanceOfJsonWebKeySpec(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function JsonWebKeySpecFromJSON(json: any): JsonWebKeySpec {
    return JsonWebKeySpecFromJSONTyped(json, false);
}

export function JsonWebKeySpecFromJSONTyped(json: any, ignoreDiscriminator: boolean): JsonWebKeySpec {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'alg': !exists(json, 'alg') ? undefined : json['alg'],
        'kty': !exists(json, 'kty') ? undefined : JsonWebKeyTypeFromJSON(json['kty']),
        'crv': !exists(json, 'crv') ? undefined : JsonWebKeyCurveNameFromJSON(json['crv']),
        'keySize': !exists(json, 'key_size') ? undefined : json['key_size'],
        'keyOps': !exists(json, 'key_ops') ? undefined : ((json['key_ops'] as Array<any>).map(JsonWebKeyOperationFromJSON)),
        'ext': !exists(json, 'ext') ? undefined : json['ext'],
    };
}

export function JsonWebKeySpecToJSON(value?: JsonWebKeySpec | null): any {
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
        'ext': value.ext,
    };
}

