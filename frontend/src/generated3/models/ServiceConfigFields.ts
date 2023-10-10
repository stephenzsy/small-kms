/* tslint:disable */
/* eslint-disable */
/**
 * Small KMS Admin API Shared
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
import type { ServiceConfigAcrInfo } from './ServiceConfigAcrInfo';
import {
    ServiceConfigAcrInfoFromJSON,
    ServiceConfigAcrInfoFromJSONTyped,
    ServiceConfigAcrInfoToJSON,
} from './ServiceConfigAcrInfo';
import type { ServiceConfigAppRoleIds } from './ServiceConfigAppRoleIds';
import {
    ServiceConfigAppRoleIdsFromJSON,
    ServiceConfigAppRoleIdsFromJSONTyped,
    ServiceConfigAppRoleIdsToJSON,
} from './ServiceConfigAppRoleIds';

/**
 * 
 * @export
 * @interface ServiceConfigFields
 */
export interface ServiceConfigFields {
    /**
     * 
     * @type {string}
     * @memberof ServiceConfigFields
     */
    azureSubscriptionId: string;
    /**
     * 
     * @type {string}
     * @memberof ServiceConfigFields
     */
    keyvaultArmResourceId: string;
    /**
     * 
     * @type {ServiceConfigAppRoleIds}
     * @memberof ServiceConfigFields
     */
    appRoleIds: ServiceConfigAppRoleIds;
    /**
     * 
     * @type {ServiceConfigAcrInfo}
     * @memberof ServiceConfigFields
     */
    azureContainerRegistry: ServiceConfigAcrInfo;
}

/**
 * Check if a given object implements the ServiceConfigFields interface.
 */
export function instanceOfServiceConfigFields(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "azureSubscriptionId" in value;
    isInstance = isInstance && "keyvaultArmResourceId" in value;
    isInstance = isInstance && "appRoleIds" in value;
    isInstance = isInstance && "azureContainerRegistry" in value;

    return isInstance;
}

export function ServiceConfigFieldsFromJSON(json: any): ServiceConfigFields {
    return ServiceConfigFieldsFromJSONTyped(json, false);
}

export function ServiceConfigFieldsFromJSONTyped(json: any, ignoreDiscriminator: boolean): ServiceConfigFields {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'azureSubscriptionId': json['azureSubscriptionId'],
        'keyvaultArmResourceId': json['keyvaultArmResourceId'],
        'appRoleIds': ServiceConfigAppRoleIdsFromJSON(json['appRoleIds']),
        'azureContainerRegistry': ServiceConfigAcrInfoFromJSON(json['azureContainerRegistry']),
    };
}

export function ServiceConfigFieldsToJSON(value?: ServiceConfigFields | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'azureSubscriptionId': value.azureSubscriptionId,
        'keyvaultArmResourceId': value.keyvaultArmResourceId,
        'appRoleIds': ServiceConfigAppRoleIdsToJSON(value.appRoleIds),
        'azureContainerRegistry': ServiceConfigAcrInfoToJSON(value.azureContainerRegistry),
    };
}

