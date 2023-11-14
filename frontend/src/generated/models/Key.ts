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
import type { JsonWebKeyType } from './JsonWebKeyType';
import {
    JsonWebKeyTypeFromJSON,
    JsonWebKeyTypeFromJSONTyped,
    JsonWebKeyTypeToJSON,
} from './JsonWebKeyType';

/**
 * 
 * @export
 * @interface Key
 */
export interface Key {
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    id: string;
    /**
     * 
     * @type {Date}
     * @memberof Key
     */
    updated: Date;
    /**
     * 
     * @type {Date}
     * @memberof Key
     */
    deleted?: Date;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    updatedBy?: string;
    /**
     * 
     * @type {number}
     * @memberof Key
     */
    iat: number;
    /**
     * 
     * @type {number}
     * @memberof Key
     */
    exp?: number;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    policy: string;
    /**
     * 
     * @type {number}
     * @memberof Key
     */
    keySize?: number;
    /**
     * 
     * @type {number}
     * @memberof Key
     */
    nbf?: number;
    /**
     * 
     * @type {JsonWebKeyType}
     * @memberof Key
     */
    kty: JsonWebKeyType;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    kid: string;
    /**
     * 
     * @type {Array<JsonWebKeyOperation>}
     * @memberof Key
     */
    keyOps?: Array<JsonWebKeyOperation>;
    /**
     * 
     * @type {JsonWebKeyCurveName}
     * @memberof Key
     */
    crv?: JsonWebKeyCurveName;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    n?: string;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    e?: string;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    x?: string;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    y?: string;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    x5u?: string;
    /**
     * Base64 encoded certificate chain
     * @type {Array<string>}
     * @memberof Key
     */
    x5c?: Array<string>;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    x5t?: string;
    /**
     * 
     * @type {string}
     * @memberof Key
     */
    x5tS256?: string;
}

/**
 * Check if a given object implements the Key interface.
 */
export function instanceOfKey(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "id" in value;
    isInstance = isInstance && "updated" in value;
    isInstance = isInstance && "iat" in value;
    isInstance = isInstance && "policy" in value;
    isInstance = isInstance && "kty" in value;
    isInstance = isInstance && "kid" in value;

    return isInstance;
}

export function KeyFromJSON(json: any): Key {
    return KeyFromJSONTyped(json, false);
}

export function KeyFromJSONTyped(json: any, ignoreDiscriminator: boolean): Key {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': json['id'],
        'updated': (new Date(json['updated'])),
        'deleted': !exists(json, 'deleted') ? undefined : (new Date(json['deleted'])),
        'updatedBy': !exists(json, 'updatedBy') ? undefined : json['updatedBy'],
        'iat': json['iat'],
        'exp': !exists(json, 'exp') ? undefined : json['exp'],
        'policy': json['policy'],
        'keySize': !exists(json, 'key_size') ? undefined : json['key_size'],
        'nbf': !exists(json, 'nbf') ? undefined : json['nbf'],
        'kty': JsonWebKeyTypeFromJSON(json['kty']),
        'kid': json['kid'],
        'keyOps': !exists(json, 'key_ops') ? undefined : ((json['key_ops'] as Array<any>).map(JsonWebKeyOperationFromJSON)),
        'crv': !exists(json, 'crv') ? undefined : JsonWebKeyCurveNameFromJSON(json['crv']),
        'n': !exists(json, 'n') ? undefined : json['n'],
        'e': !exists(json, 'e') ? undefined : json['e'],
        'x': !exists(json, 'x') ? undefined : json['x'],
        'y': !exists(json, 'y') ? undefined : json['y'],
        'x5u': !exists(json, 'x5u') ? undefined : json['x5u'],
        'x5c': !exists(json, 'x5c') ? undefined : json['x5c'],
        'x5t': !exists(json, 'x5t') ? undefined : json['x5t'],
        'x5tS256': !exists(json, 'x5t#S256') ? undefined : json['x5t#S256'],
    };
}

export function KeyToJSON(value?: Key | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'updated': (value.updated.toISOString()),
        'deleted': value.deleted === undefined ? undefined : (value.deleted.toISOString()),
        'updatedBy': value.updatedBy,
        'iat': value.iat,
        'exp': value.exp,
        'policy': value.policy,
        'key_size': value.keySize,
        'nbf': value.nbf,
        'kty': JsonWebKeyTypeToJSON(value.kty),
        'kid': value.kid,
        'key_ops': value.keyOps === undefined ? undefined : ((value.keyOps as Array<any>).map(JsonWebKeyOperationToJSON)),
        'crv': JsonWebKeyCurveNameToJSON(value.crv),
        'n': value.n,
        'e': value.e,
        'x': value.x,
        'y': value.y,
        'x5u': value.x5u,
        'x5c': value.x5c,
        'x5t': value.x5t,
        'x5t#S256': value.x5tS256,
    };
}
