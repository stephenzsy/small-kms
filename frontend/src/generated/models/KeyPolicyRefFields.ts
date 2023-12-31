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
/**
 * 
 * @export
 * @interface KeyPolicyRefFields
 */
export interface KeyPolicyRefFields {
    /**
     * 
     * @type {string}
     * @memberof KeyPolicyRefFields
     */
    displayName: string;
}

/**
 * Check if a given object implements the KeyPolicyRefFields interface.
 */
export function instanceOfKeyPolicyRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function KeyPolicyRefFieldsFromJSON(json: any): KeyPolicyRefFields {
    return KeyPolicyRefFieldsFromJSONTyped(json, false);
}

export function KeyPolicyRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): KeyPolicyRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': json['displayName'],
    };
}

export function KeyPolicyRefFieldsToJSON(value?: KeyPolicyRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'displayName': value.displayName,
    };
}

