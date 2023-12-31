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
 * @interface ProfileParameters
 */
export interface ProfileParameters {
    /**
     * 
     * @type {string}
     * @memberof ProfileParameters
     */
    displayName?: string;
}

/**
 * Check if a given object implements the ProfileParameters interface.
 */
export function instanceOfProfileParameters(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function ProfileParametersFromJSON(json: any): ProfileParameters {
    return ProfileParametersFromJSONTyped(json, false);
}

export function ProfileParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): ProfileParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': !exists(json, 'displayName') ? undefined : json['displayName'],
    };
}

export function ProfileParametersToJSON(value?: ProfileParameters | null): any {
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

