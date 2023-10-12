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

import {
    CreateManagedApplicationProfileRequest,
    instanceOfCreateManagedApplicationProfileRequest,
    CreateManagedApplicationProfileRequestFromJSON,
    CreateManagedApplicationProfileRequestFromJSONTyped,
    CreateManagedApplicationProfileRequestToJSON,
} from './CreateManagedApplicationProfileRequest';

/**
 * @type CreateProfileRequest
 * 
 * @export
 */
export type CreateProfileRequest = { type: 'managed-application' } & CreateManagedApplicationProfileRequest;

export function CreateProfileRequestFromJSON(json: any): CreateProfileRequest {
    return CreateProfileRequestFromJSONTyped(json, false);
}

export function CreateProfileRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): CreateProfileRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    switch (json['type']) {
        case 'managed-application':
            return {...CreateManagedApplicationProfileRequestFromJSONTyped(json, true), type: 'managed-application'};
        default:
            throw new Error(`No variant of CreateProfileRequest exists with 'type=${json['type']}'`);
    }
}

export function CreateProfileRequestToJSON(value?: CreateProfileRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    switch (value['type']) {
        case 'managed-application':
            return CreateManagedApplicationProfileRequestToJSON(value);
        default:
            throw new Error(`No variant of CreateProfileRequest exists with 'type=${value['type']}'`);
    }

}

