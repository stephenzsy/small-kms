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
/**
 * 
 * @export
 * @interface RequestHeaderEntry
 */
export interface RequestHeaderEntry {
    /**
     * 
     * @type {string}
     * @memberof RequestHeaderEntry
     */
    key: string;
    /**
     * 
     * @type {Array<string>}
     * @memberof RequestHeaderEntry
     */
    value: Array<string>;
}

/**
 * Check if a given object implements the RequestHeaderEntry interface.
 */
export function instanceOfRequestHeaderEntry(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "key" in value;
    isInstance = isInstance && "value" in value;

    return isInstance;
}

export function RequestHeaderEntryFromJSON(json: any): RequestHeaderEntry {
    return RequestHeaderEntryFromJSONTyped(json, false);
}

export function RequestHeaderEntryFromJSONTyped(json: any, ignoreDiscriminator: boolean): RequestHeaderEntry {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'key': json['key'],
        'value': json['value'],
    };
}

export function RequestHeaderEntryToJSON(value?: RequestHeaderEntry | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'key': value.key,
        'value': value.value,
    };
}

