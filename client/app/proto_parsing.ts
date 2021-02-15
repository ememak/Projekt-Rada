import { EnvelopeToSign, PollSchema, RSASignature, VoteRequest } from "Projekt_Rada/query/query_pb";

export function QAListToSchema(questionsList: PollSchema.QA.AsObject[]) {
  let schema = new PollSchema();
  for (let qa of questionsList){
    const QA = new PollSchema.QA();
    QA.setQuestion(qa.question);
    QA.setType(qa.type as 0 | 1 | 2);
    QA.setOptionsList(qa.optionsList);
    QA.setAnswersList(qa.answersList);
    schema.addQuestions(QA);
  }
  return schema
}

export function toEnvelope(envelope: string, pollid: number, token: string) {
  let request: EnvelopeToSign = new EnvelopeToSign();
  request.setEnvelope(hexToBase64(envelope))
  request.setPollid(pollid)
  request.setToken(token)
  return request
}

export function toRSASignature(ballot, sign: string) {
  let signature: RSASignature = new RSASignature();

  signature.setBallot(hexToBase64(ballot));
  signature.setSign(hexToBase64(sign));
  return signature;
}

export function toVoteRequest(pollid: number, schema: PollSchema, signature: RSASignature) {
  let request: VoteRequest = new VoteRequest();

  request.setPollid(pollid)
  request.setAnswers(schema)
  request.setSign(signature)
  return request
}

function hexToBase64(hexstring) {
    return btoa(hexstring.match(/\w{2}/g).map(function(a) {
        return String.fromCharCode(parseInt(a, 16));
    }).join(""));
}
