import { GetCertificateV1Request } from "../generated";
import * as runtime from "../generated/runtime";

/**
 *
 */
export class CertsDownloadApi extends runtime.BaseAPI {
  /**
   * Get certificate
   */
  async getCertificateV1Raw(
    requestParameters: GetCertificateV1Request,
    initOverrides?: RequestInit | runtime.InitOverrideFunction
  ): Promise<Blob> {
    if (
      requestParameters.namespaceId === null ||
      requestParameters.namespaceId === undefined
    ) {
      throw new runtime.RequiredError(
        "namespaceId",
        "Required parameter requestParameters.namespaceId was null or undefined when calling getCertificateV1."
      );
    }

    if (requestParameters.id === null || requestParameters.id === undefined) {
      throw new runtime.RequiredError(
        "id",
        "Required parameter requestParameters.id was null or undefined when calling getCertificateV1."
      );
    }

    const queryParameters: any = {};

    const headerParameters: runtime.HTTPHeaders = {};

    if (
      requestParameters.accept !== undefined &&
      requestParameters.accept !== null
    ) {
      headerParameters["Accept"] = String(requestParameters.accept);
    }

    if (this.configuration && this.configuration.accessToken) {
      const token = this.configuration.accessToken;
      const tokenString = await token("BearerAuth", []);

      if (tokenString) {
        headerParameters["Authorization"] = `Bearer ${tokenString}`;
      }
    }
    const response = await this.request(
      {
        path: `/v1/{namespaceId}/certificate/{id}`
          .replace(
            `{${"namespaceId"}}`,
            encodeURIComponent(String(requestParameters.namespaceId))
          )
          .replace(
            `{${"id"}}`,
            encodeURIComponent(String(requestParameters.id))
          ),
        method: "GET",
        headers: headerParameters,
        query: queryParameters,
      },
      initOverrides
    );
    return response.blob();
  }

  /**
   * Get certificate
   */
  getCertificateV1(
    requestParameters: GetCertificateV1Request,
    initOverrides?: RequestInit | runtime.InitOverrideFunction
  ): Promise<Blob> {
    return this.getCertificateV1Raw(requestParameters, initOverrides);
  }
}
