import { useRequest } from "ahooks";
import { Card, CardSection, CartTitle } from "../components/Card";
import { useAuthedClient } from "../utils/useCertsApi3";
import { AdminApi } from "../generated3";

export default function ServicePage() {
  const adminApi = useAuthedClient(AdminApi);
  const { data: serviceConfig } = useRequest(() => {
    return adminApi.getServiceConfig();
  });
  return (
    <Card>
      <CartTitle>Service configuration</CartTitle>
      <CardSection>
        <pre>{JSON.stringify(serviceConfig, null, 2)}</pre>
      </CardSection>
    </Card>
  );
}
