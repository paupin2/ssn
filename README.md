Package `sessn` validades, normalizes and generates [Swedish SSNs](https://en.wikipedia.org/wiki/Personal_identity_number_(Sweden)),
or _Personnummer_. It has no external dependencies and supports regular,
coordination and interim numbers, and is written with clarity
and ease of use in mind.

Use `sessn.Normalize` to check/normalize numbers before storing them, since it's
[recommended](https://www.riksdagen.se/sv/dokument-och-lagar/dokument/svensk-forfattningssamling/folkbokforingslag-1991481_sfs-1991-481/#P18)
to store the 12-digit version, to avoid ambiguity.

The following formats (with whitespace trimmed) are recognized
(`C`:century, `Y`:year, `M`:month, `D`:day, `I`:interim letter, `B`:birth number, `K`:check digit):

- Regular SSN:
  - "CCYYMMDD-BBBK"
  - "CCYYMMDD+BBBK"
  - "CCYYMMDDBBBK"
  - "YYMMDD-BBBK"
  - "YYMMDD+BBBK"
  - "YYMMDDBBBK"
- Interim number:
  - "CCYYMMDD-IBBK"
  - "CCYYMMDD+IBBK"
  - "CCYYMMDDIBBK"
  - "YYMMDD-IBBK"
  - "YYMMDD+IBBK"
  - "YYMMDDIBBK"

Test data from [here](https://www7.skatteverket.se/portal/apier-och-oppna-data/utvecklarportalen/oppetdata/Test%C2%AD%C2%ADpersonnummer).
