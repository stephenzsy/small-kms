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
 * @interface AzureRoleAssignment
 */
export interface AzureRoleAssignment {
    /**
     * 
     * @type {string}
     * @memberof AzureRoleAssignment
     */
    id?: string;
    /**
     * 
     * @type {string}
     * @memberof AzureRoleAssignment
     */
    name?: string;
    /**
     * 
     * @type {string}
     * @memberof AzureRoleAssignment
     */
    roleDefinitionId?: string;
    /**
     * 
     * @type {string}
     * @memberof AzureRoleAssignment
     */
    principalId?: string;
}

/**
 * Check if a given object implements the AzureRoleAssignment interface.
 */
export function instanceOfAzureRoleAssignment(value: object): boolean {
    let isInstance = true;

    return isInstance;
}

export function AzureRoleAssignmentFromJSON(json: any): AzureRoleAssignment {
    return AzureRoleAssignmentFromJSONTyped(json, false);
}

export function AzureRoleAssignmentFromJSONTyped(json: any, ignoreDiscriminator: boolean): AzureRoleAssignment {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'id': !exists(json, 'id') ? undefined : json['id'],
        'name': !exists(json, 'name') ? undefined : json['name'],
        'roleDefinitionId': !exists(json, 'roleDefinitionId') ? undefined : json['roleDefinitionId'],
        'principalId': !exists(json, 'principalId') ? undefined : json['principalId'],
    };
}

export function AzureRoleAssignmentToJSON(value?: AzureRoleAssignment | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'id': value.id,
        'name': value.name,
        'roleDefinitionId': value.roleDefinitionId,
        'principalId': value.principalId,
    };
}

