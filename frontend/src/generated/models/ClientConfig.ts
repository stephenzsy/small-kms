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
 * @interface ClientConfig
 */
export interface ClientConfig {
    /**
     * 
     * @type {string}
     * @memberof ClientConfig
     */
    name: string;
    /**
     * 
     * @type {string}
     * @memberof ClientConfig
     */
    ipaddr?: string;
    /**
     * 
     * @type {string}
     * @memberof ClientConfig
     */
    secret?: string;
}

/**
 * Check if a given object implements the ClientConfig interface.
 */
export function instanceOfClientConfig(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "name" in value;

    return isInstance;
}

export function ClientConfigFromJSON(json: any): ClientConfig {
    return ClientConfigFromJSONTyped(json, false);
}

export function ClientConfigFromJSONTyped(json: any, ignoreDiscriminator: boolean): ClientConfig {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'name': json['name'],
        'ipaddr': !exists(json, 'ipaddr') ? undefined : json['ipaddr'],
        'secret': !exists(json, 'secret') ? undefined : json['secret'],
    };
}

export function ClientConfigToJSON(value?: ClientConfig | null): any {
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
    };
}

