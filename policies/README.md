# AWS policies for the Parameter Store
These example policies should be copied and edited to suit your preferences, then applied as policies in AWS Identity and Access Management (IAM).

- [AdminParameterStore](AdminParameterStore.json) defines a full access policy to the Parameter Store.
- [ReadParameterStore](ReadParameterStore.json) defines read-only permissions to keys starting with "/secret".
- [WriteParameterStore](WriteParameterStore.json) defines read/write access to all keys starting with "/secret", and the ability to create new keys with the same pattern.

NOTE: The ssm:DescribeParameters permission (for listing contents) only accepts "*" as its resource, so any users with access at all will also be able to see the names of keys. Plan accordingly.

## Usage
- Log on to your AWS console (https://console.aws.amazon.com/)
- Go to My Security Credentials from the account dropdown
- Add each of the three policies under Policies.
- Add a user group for each and add users as needed.

It may take a couple of minutes before policies activate.
