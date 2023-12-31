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
 * @interface PullImageRequest
 */
export interface PullImageRequest {
    /**
     * 
     * @type {string}
     * @memberof PullImageRequest
     */
    imageRepo?: string;
    /**
     * 
     * @type {string}
     * @memberof PullImageRequest
     */
    imageTag: string;
    /**
     * 
     * @type {boolean}
     * @memberof PullImageRequest
     */
    includeLatestTag?: boolean;
}

/**
 * Check if a given object implements the PullImageRequest interface.
 */
export function instanceOfPullImageRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "imageTag" in value;

    return isInstance;
}

export function PullImageRequestFromJSON(json: any): PullImageRequest {
    return PullImageRequestFromJSONTyped(json, false);
}

export function PullImageRequestFromJSONTyped(json: any, ignoreDiscriminator: boolean): PullImageRequest {
    if ((json === undefined) || (json === null)) {
        return json;
    }
    return {
        
        'imageRepo': !exists(json, 'imageRepo') ? undefined : json['imageRepo'],
        'imageTag': json['imageTag'],
        'includeLatestTag': !exists(json, 'includeLatestTag') ? undefined : json['includeLatestTag'],
    };
}

export function PullImageRequestToJSON(value?: PullImageRequest | null): any {
    if (value === undefined) {
        return undefined;
    }
    if (value === null) {
        return null;
    }
    return {
        
        'imageRepo': value.imageRepo,
        'imageTag': value.imageTag,
        'includeLatestTag': value.includeLatestTag,
    };
}

