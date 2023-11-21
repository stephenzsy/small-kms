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
import type { KeyStatus } from './KeyStatus';
import {
    KeyStatusFromJSON,
    KeyStatusFromJSONTyped,
    KeyStatusToJSON,
} from './KeyStatus';

/**
 * 
 * @export
 * @interface KeyRefFields
 */
export interface KeyRefFields {
    /**
     * 
     * @type {KeyStatus}
     * @memberof KeyRefFields
     */
    status: KeyStatus;
    /**
     * 
     * @type {number}
     * @memberof KeyRefFields
     */
    iat: number;
    /**
     * 
     * @type {number}
     * @memberof KeyRefFields
     */
    exp?: number;
    /**
     * 
     * @type {string}
     * @memberof KeyRefFields
     */
    policyIdentifier: string;
}

/**
 * Check if a given object implements the KeyRefFields interface.
 */
export function instanceOfKeyRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "status" in value;
    isInstance = isInstance && "iat" in value;
    isInstance = isInstance && "policyIdentifier" in value;

    return isInstance;
}

export function KeyRefFieldsFromJSON(json: any): KeyRefFields {
    return KeyRefFieldsFromJSONTyped(json, false);
}

export function KeyRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): KeyRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'status': KeyStatusFromJSON(json['status']),
        'iat': json['iat'],
        'exp': !exists(json, 'exp') ? undefined : json['exp'],
        'policyIdentifier': json['policyIdentifier'],
    };
}

export function KeyRefFieldsToJSON(value?: KeyRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'status': KeyStatusToJSON(value.status),
        'iat': value.iat,
        'exp': value.exp,
        'policyIdentifier': value.policyIdentifier,
    };
}

