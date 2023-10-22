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
 * @interface ProfileRefFields1
 */
export interface ProfileRefFields1 {
    /**
     * 
     * @type {string}
     * @memberof ProfileRefFields1
     */
    displayName: string;
}

/**
 * Check if a given object implements the ProfileRefFields1 interface.
 */
export function instanceOfProfileRefFields1(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function ProfileRefFields1FromJSON(json: any): ProfileRefFields1 {
    return ProfileRefFields1FromJSONTyped(json, false);
}

export function ProfileRefFields1FromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileRefFields1 {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': json['displayName'],
    };
}

export function ProfileRefFields1ToJSON(value?: ProfileRefFields1 | null): any {
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

