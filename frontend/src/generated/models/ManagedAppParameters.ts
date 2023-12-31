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
 * @interface ManagedAppParameters
 */
export interface ManagedAppParameters {
    /**
     * 
     * @type {string}
     * @memberof ManagedAppParameters
     */
    displayName: string;
    /**
     * 
     * @type {boolean}
     * @memberof ManagedAppParameters
     */
    skipServicePrincipalCreation?: boolean;
}

/**
 * Check if a given object implements the ManagedAppParameters interface.
 */
export function instanceOfManagedAppParameters(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "displayName" in value;

    return isInstance;
}

export function ManagedAppParametersFromJSON(json: any): ManagedAppParameters {
    return ManagedAppParametersFromJSONTyped(json, false);
}

export function ManagedAppParametersFromJSONTyped(json: any, ignoreDiscriminator: boolean): ManagedAppParameters {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'displayName': json['displayName'],
        'skipServicePrincipalCreation': !exists(json, 'skipServicePrincipalCreation') ? undefined : json['skipServicePrincipalCreation'],
    };
}

export function ManagedAppParametersToJSON(value?: ManagedAppParameters | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'displayName': value.displayName,
        'skipServicePrincipalCreation': value.skipServicePrincipalCreation,
    };
}

