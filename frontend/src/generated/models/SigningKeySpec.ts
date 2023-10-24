/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 0.1.1
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
import type { JsonWebKeySignatureAlgorithm } from './JsonWebKeySignatureAlgorithm';
import {
    JsonWebKeySignatureAlgorithmFromJSON,
    JsonWebKeySignatureAlgorithmFromJSONTyped,
    JsonWebKeySignatureAlgorithmToJSON,
} from './JsonWebKeySignatureAlgorithm';
import type { JsonWebKeyType } from './JsonWebKeyType';
import {
    JsonWebKeyTypeFromJSON,
    JsonWebKeyTypeFromJSONTyped,
    JsonWebKeyTypeToJSON,
} from './JsonWebKeyType';

/**
 * 
 * @export
 * @interface SigningKeySpec
 */
export interface SigningKeySpec {
    /**
     * 
     * @type {JsonWebKeyType}
     * @memberof SigningKeySpec
     */
    kty: JsonWebKeyType;
    /**
     * 
     * @type {string}
     * @memberof SigningKeySpec
     */
    kid?: string;
    /**
     * 
     * @type {JsonWebKeyCurveName}
     * @memberof SigningKeySpec
     */
    crv?: JsonWebKeyCurveName;
    /**
     * 
     * @type {number}
     * @memberof SigningKeySpec
     */
    keySize?: number;
    /**
     * 
     * @type {Array<JsonWebKeyOperation>}
     * @memberof SigningKeySpec
     */
    keyOps: Array<JsonWebKeyOperation>;
    /**
     * 
     * @type {JsonWebKeySignatureAlgorithm}
     * @memberof SigningKeySpec
     */
    alg?: JsonWebKeySignatureAlgorithm;
    /**
     * Base64 encoded certificate chain
     * @type {Array<string>}
     * @memberof SigningKeySpec
     */
    x5c?: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof SigningKeySpec
     */
    x5t?: string;
    /**
     * 
     * @type {string}
     * @memberof SigningKeySpec
     */
    x5tS256?: string;
}

/**
 * Check if a given object implements the SigningKeySpec interface.
 */
export function instanceOfSigningKeySpec(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "kty" in value;
    isInstance = isInstance && "keyOps" in value;

    return isInstance;
}

export function SigningKeySpecFromJSON(json: any): SigningKeySpec {
    return SigningKeySpecFromJSONTyped(json, false);
}

export function SigningKeySpecFromJSONTyped(json: any, ignoreDiscriminator: boolean): SigningKeySpec {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'kty': JsonWebKeyTypeFromJSON(json['kty']),
        'kid': !exists(json, 'kid') ? undefined : json['kid'],
        'crv': !exists(json, 'crv') ? undefined : JsonWebKeyCurveNameFromJSON(json['crv']),
        'keySize': !exists(json, 'key_size') ? undefined : json['key_size'],
        'keyOps': ((json['key_ops'] as Array<any>).map(JsonWebKeyOperationFromJSON)),
        'alg': !exists(json, 'alg') ? undefined : JsonWebKeySignatureAlgorithmFromJSON(json['alg']),
        'x5c': !exists(json, 'x5c') ? undefined : json['x5c'],
        'x5t': !exists(json, 'x5t') ? undefined : json['x5t'],
        'x5tS256': !exists(json, 'x5t#S256') ? undefined : json['x5t#S256'],
    };
}

export function SigningKeySpecToJSON(value?: SigningKeySpec | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'kty': JsonWebKeyTypeToJSON(value.kty),
        'kid': value.kid,
        'crv': JsonWebKeyCurveNameToJSON(value.crv),
        'key_size': value.keySize,
        'key_ops': ((value.keyOps as Array<any>).map(JsonWebKeyOperationToJSON)),
        'alg': JsonWebKeySignatureAlgorithmToJSON(value.alg),
        'x5c': value.x5c,
        'x5t': value.x5t,
        'x5t#S256': value.x5tS256,
    };
}

