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
import type { JwkAlg } from './JwkAlg';
import {
    JwkAlgFromJSON,
    JwkAlgFromJSONTyped,
    JwkAlgToJSON,
} from './JwkAlg';
import type { JwkKeySize } from './JwkKeySize';
import {
    JwkKeySizeFromJSON,
    JwkKeySizeFromJSONTyped,
    JwkKeySizeToJSON,
} from './JwkKeySize';
import type { KeyType } from './KeyType';
import {
    KeyTypeFromJSON,
    KeyTypeFromJSONTyped,
    KeyTypeToJSON,
} from './KeyType';

/**
 * Partial implementation of JSON Web Key (RFC 7517) with additional fields
 * @export
 * @interface JwkKeyProperties
 */
export interface JwkKeyProperties {
    /**
     * 
     * @type {JwkAlg}
     * @memberof JwkKeyProperties
     */
    alg?: JwkAlg;
    /**
     * 
     * @type {KeyType}
     * @memberof JwkKeyProperties
     */
    kty: KeyType;
    /**
     * 
     * @type {JwkKeySize}
     * @memberof JwkKeyProperties
     */
    keySize?: JwkKeySize;
    /**
     * 
     * @type {string}
     * @memberof JwkKeyProperties
     */
    n?: string;
    /**
     * 
     * @type {string}
     * @memberof JwkKeyProperties
     */
    e?: string;
    /**
     * 
     * @type {CurveName}
     * @memberof JwkKeyProperties
     */
    crv?: CurveName;
    /**
     * 
     * @type {string}
     * @memberof JwkKeyProperties
     */
    x?: string;
    /**
     * 
     * @type {string}
     * @memberof JwkKeyProperties
     */
    y?: string;
}

/**
 * Check if a given object implements the JwkKeyProperties interface.
 */
export function instanceOfJwkKeyProperties(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "kty" in value;

    return isInstance;
}

export function JwkKeyPropertiesFromJSON(json: any): JwkKeyProperties {
    return JwkKeyPropertiesFromJSONTyped(json, false);
}

export function JwkKeyPropertiesFromJSONTyped(json: any, ignoreDiscriminator: boolean): JwkKeyProperties {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'alg': !exists(json, 'alg') ? undefined : JwkAlgFromJSON(json['alg']),
        'kty': KeyTypeFromJSON(json['kty']),
        'keySize': !exists(json, 'key_size') ? undefined : JwkKeySizeFromJSON(json['key_size']),
        'n': !exists(json, 'n') ? undefined : json['n'],
        'e': !exists(json, 'e') ? undefined : json['e'],
        'crv': !exists(json, 'crv') ? undefined : CurveNameFromJSON(json['crv']),
        'x': !exists(json, 'x') ? undefined : json['x'],
        'y': !exists(json, 'y') ? undefined : json['y'],
    };
}

export function JwkKeyPropertiesToJSON(value?: JwkKeyProperties | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'alg': JwkAlgToJSON(value.alg),
        'kty': KeyTypeToJSON(value.kty),
        'key_size': JwkKeySizeToJSON(value.keySize),
        'n': value.n,
        'e': value.e,
        'crv': CurveNameToJSON(value.crv),
        'x': value.x,
        'y': value.y,
    };
}
