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
import type { JsonWebKeySpec } from './JsonWebKeySpec';
import {
    JsonWebKeySpecFromJSON,
    JsonWebKeySpecFromJSONTyped,
    JsonWebKeySpecToJSON,
} from './JsonWebKeySpec';

/**
 * 
 * @export
 * @interface CreateKeyPolicyRequest
 */
export interface CreateKeyPolicyRequest {
    /**
     * 
     * @type {string}
     * @memberof CreateKeyPolicyRequest
     */
    displayName?: string;
    /**
     * 
     * @type {JsonWebKeySpec}
     * @memberof CreateKeyPolicyRequest
     */
    keySpec?: JsonWebKeySpec;
    /**
     * 
     * @type {boolean}
     * @memberof CreateKeyPolicyRequest
     */
    exportable?: boolean;
    /**
     * 
     * @type {string}
     * @memberof CreateKeyPolicyRequest
     */
    expiryTime?: string;
}

/**
 * Check if a given object implements the CreateKeyPolicyRequest interface.
 */
export function instanceOfCreateKeyPolicyRequest(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function CreateKeyPolicyRequestFromJSON(json: any): CreateKeyPolicyRequest {
    return CreateKeyPolicyRequestFromJSONTyped(json, false);
}

export function CreateKeyPolicyRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateKeyPolicyRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
        'keySpec': !exists(json, 'keySpec') ? undefined : JsonWebKeySpecFromJSON(json['keySpec']),
        'exportable': !exists(json, 'exportable') ? undefined : json['exportable'],
        'expiryTime': !exists(json, 'expiryTime') ? undefined : json['expiryTime'],
    };
}

export function CreateKeyPolicyRequestToJSON(value?: CreateKeyPolicyRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'displayName': value.displayName,
        'keySpec': JsonWebKeySpecToJSON(value.keySpec),
        'exportable': value.exportable,
        'expiryTime': value.expiryTime,
    };
}

