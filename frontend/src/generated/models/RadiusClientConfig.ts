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
 * @interface RadiusClientConfig
 */
export interface RadiusClientConfig {
    /**
     * 
     * @type {string}
     * @memberof RadiusClientConfig
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof RadiusClientConfig
     */
    ipaddr?: string;
    /**
     * 
     * @type {string}
     * @memberof RadiusClientConfig
     */
    secret?: string;
    /**
     * 
     * @type {string}
     * @memberof RadiusClientConfig
     */
    secretId?: string;
    /**
     * 
     * @type {string}
     * @memberof RadiusClientConfig
     */
    secretPolicyId?: string;
}

/**
 * Check if a given object implements the RadiusClientConfig interface.
 */
export function instanceOfRadiusClientConfig(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "name" in value;

    return isInstance;
}

export function RadiusClientConfigFromJSON(json: any): RadiusClientConfig {
    return RadiusClientConfigFromJSONTyped(json, false);
}

export function RadiusClientConfigFromJSONTyped(json: any, ignoreDiscriminator: boolean): RadiusClientConfig {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'name': json['name'],
        'ipaddr': !exists(json, 'ipaddr') ? undefined : json['ipaddr'],
        'secret': !exists(json, 'secret') ? undefined : json['secret'],
        'secretId': !exists(json, 'secretId') ? undefined : json['secretId'],
        'secretPolicyId': !exists(json, 'secretPolicyId') ? undefined : json['secretPolicyId'],
    };
}

export function RadiusClientConfigToJSON(value?: RadiusClientConfig | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'name': value.name,
        'ipaddr': value.ipaddr,
        'secret': value.secret,
        'secretId': value.secretId,
        'secretPolicyId': value.secretPolicyId,
    };
}

