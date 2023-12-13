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
 * @interface PullImageRequest
 */
export interface PullImageRequest {
    /**
     * 
     * @type {string}
     * @memberof PullImageRequest
     */
    imageRepo: string;
    /**
     * 
     * @type {string}
     * @memberof PullImageRequest
     */
    imageTag: string;
}

/**
 * Check if a given object implements the PullImageRequest interface.
 */
export function instanceOfPullImageRequest(value: object): boolean {
    let isInstance = true;
    isInstance = isInstance && "imageRepo" in value;
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
        
        'imageRepo': json['imageRepo'],
        'imageTag': json['imageTag'],
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
    };
}

