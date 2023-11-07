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
 * @interface SecretRefFields
 */
export interface SecretRefFields {
    /**
     * 
     * @type {string}
     * @memberof SecretRefFields
     */
    version: string;
    /**
     * 
     * @type {string}
     * @memberof SecretRefFields
     */
    policyId: string;
}

/**
 * Check if a given object implements the SecretRefFields interface.
 */
export function instanceOfSecretRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "version" in value;
    isInstance = isInstance && "policyId" in value;

    return isInstance;
}

export function SecretRefFieldsFromJSON(json: any): SecretRefFields {
    return SecretRefFieldsFromJSONTyped(json, false);
}

export function SecretRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): SecretRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'version': json['version'],
        'policyId': json['policyId'],
    };
}

export function SecretRefFieldsToJSON(value?: SecretRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'version': value.version,
        'policyId': value.policyId,
    };
}

