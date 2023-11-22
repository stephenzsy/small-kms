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
 * @interface LinkRefFields
 */
export interface LinkRefFields {
    /**
     * 
     * @type {string}
     * @memberof LinkRefFields
     */
    linkTo: string;
}

/**
 * Check if a given object implements the LinkRefFields interface.
 */
export function instanceOfLinkRefFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "linkTo" in value;

    return isInstance;
}

export function LinkRefFieldsFromJSON(json: any): LinkRefFields {
    return LinkRefFieldsFromJSONTyped(json, false);
}

export function LinkRefFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): LinkRefFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'linkTo': json['linkTo'],
    };
}

export function LinkRefFieldsToJSON(value?: LinkRefFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'linkTo': value.linkTo,
    };
}
