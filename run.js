#!/usr/bin/env node
/* eslint-disable no-await-in-loop */

const request = require('superagent');
const bluebird = require('bluebird');
const _ = require('lodash');

const organization = process.argv[2];
const token = process.env.GITHUB_TOKEN;

const OUR_LABEL_TEXT = `answering: reported by ${organization}`;
const ANSWERED_LABEL_TEXT = 'answering: answered';
const NOT_ANSWERED_LABEL_TEXT = 'answering: not answered';

const answeringLabels = [
  { name: OUR_LABEL_TEXT, color: 'a0a000' },
  { name: ANSWERED_LABEL_TEXT, color: '00a000' },
  { name: NOT_ANSWERED_LABEL_TEXT, color: 'a00000' },
];

const botsList = [
  'issuehunt-app',
];

const findRepos = async () => {
  const repoNames = [];
  for (let page = 1; ; page += 1) {
    const data = await request(`https://api.github.com/orgs/${organization}/repos`)
      .query({ page })
      .set('Authorization', `token ${token}`)
      .set('user-agent', 'script')
      .then(({ body }) => body)
      .catch((error) => {
        console.error(_.get(error, 'response.error', error));
        throw error;
      });
    if (data.length < 1) {
      break;
    }
    repoNames.push(...data.filter((repo) => !repo.archived).map((repo) => repo.name));
  }

  return repoNames.sort();
};

const findLabels = (repo) => request(`https://api.github.com/repos/${organization}/${repo}/labels`)
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .then(({ body }) => body)
  .catch((error) => {
    console.error(_.get(error, 'response.error', error));
    throw error;
  });

const createLabel = (repo, label) => request.post(`https://api.github.com/repos/${organization}/${repo}/labels`)
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .send(label)
  .then(({ body }) => body)
  .catch((error) => {
    console.error(_.get(error, 'response.error', error));
    throw error;
  });

const deleteLabel = (repo, labelName) => request.delete(`https://api.github.com/repos/${organization}/${repo}/labels/${labelName}`)
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .then(({ body }) => body)
  .catch((error) => {
    console.error(_.get(error, 'response.error', error));
    throw error;
  });

const addLabel = (issueUrl, labelName) => request.post(`${issueUrl.replace('https://github.com', 'https://api.github.com/repos')}/labels`)
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .send({ labels: [labelName] })
  .then(({ body }) => body)
  .catch((error) => {
    console.error(_.get(error, 'response.error', error));
    throw error;
  });

const removeLabel = (issueUrl, labelName) => request.delete(`${issueUrl.replace('https://github.com', 'https://api.github.com/repos')}/labels/${labelName}`)
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .then(({ body }) => body)
  .catch((error) => {
    console.error(_.get(error, 'response.error', error));
    throw error;
  });

const graphql = ({ query, variables }) => request.post('https://api.github.com/graphql')
  .set('Authorization', `token ${token}`)
  .set('user-agent', 'script')
  .send({ query, variables })
  .then(({ body }) => {
    if (body.errors) {
      console.error(body.errors);
      throw new Error();
    }

    return body;
  });

const findIssues = async (repo) => {
  const result = [];
  let cursor = null;
  while (true) {
    console.log({ cursor });
    const info = await graphql({
      query: `
  query ($organization: String!, $repo: String!, $cursor: String) {
  repository(owner: $organization, name: $repo) {
    issues(first:20, after: $cursor, states:OPEN) {
      edges {
        cursor
        node {
          title
          url
          number
          authorAssociation
          comments(last:100) {
            edges {
              node {
                bodyText
                authorAssociation
                author {
                  login
                }
              }
            }
          }
        }
      }
    }
  }
  }`,
      variables: { organization, repo, cursor },
    });
    const issues = info.data.repository.issues.edges;
    if (!issues.length) {
      break;
    }
    result.push(...issues);
    cursor = _.last(issues).cursor;
  }
  return result;
};

const updateLabels = async () => {
  const repos = await findRepos();
  console.log('found GitHub repos:');
  repos.forEach((name) => console.log(name));
  console.log();

  const allInfo = await bluebird.mapSeries(repos, async (repo) => {
    console.log({ organization, repo });
    const allLabels = await findLabels(repo);
    const labelsToRemove = answeringLabels.filter((label) => allLabels.some((anyLabel) => anyLabel.name === label.name && anyLabel.color !== label.color));
    const labelsToAdd = answeringLabels.filter(
      (label) => !allLabels.some((anyLabel) => anyLabel.name === label.name)
        || allLabels.some((anyLabel) => anyLabel.name === label.name && anyLabel.color !== label.color),
    );
    console.log({ labelsToRemove, labelsToAdd });
    await Promise.all(labelsToRemove.map((label) => deleteLabel(repo, label.name)));
    await Promise.all(labelsToAdd.map((label) => createLabel(repo, label)));
    const issues = await findIssues(repo);
    console.log(issues);
    return issues;
  });
  const allIssues = _.flatten(allInfo).map((edge) => edge.node);
  const ourIssues = allIssues.filter((issue) => issue.authorAssociation === 'MEMBER');
  const notOurIssues = allIssues.filter((issue) => issue.authorAssociation !== 'MEMBER');
  const notOurIssuesWithNoComments = notOurIssues.filter((issue) => issue.comments.edges.length === 0);
  const notOurIssuesWithComments = notOurIssues.filter((issue) => issue.comments.edges.length !== 0);
  const answeredIssues = notOurIssuesWithComments.filter(
    (issue) => {
      const lastComment = _.findLast(issue.comments.edges, (edge) => !botsList.includes(edge.node.author.login));
      if (!lastComment) {
        return false;
      }
      return lastComment.node.authorAssociation === 'MEMBER';
    },
  );
  const notAnsweredIssues = [
    ...notOurIssuesWithNoComments,
    ...notOurIssuesWithComments.filter((issue) => !answeredIssues.some((answeredIssue) => answeredIssue.url === issue.url)),
  ];
  console.log('answeredIssues', answeredIssues);
  console.log('notAnsweredIssues', notAnsweredIssues);

  await bluebird.mapSeries(allIssues, (issue) => removeLabel(issue.url, OUR_LABEL_TEXT).catch(() => {}));
  await bluebird.mapSeries(allIssues, (issue) => removeLabel(issue.url, ANSWERED_LABEL_TEXT).catch(() => {}));
  await bluebird.mapSeries(allIssues, (issue) => removeLabel(issue.url, NOT_ANSWERED_LABEL_TEXT).catch(() => {}));
  await bluebird.mapSeries(ourIssues, (issue) => addLabel(issue.url, OUR_LABEL_TEXT));
  await bluebird.mapSeries(answeredIssues, (issue) => addLabel(issue.url, ANSWERED_LABEL_TEXT));
  await bluebird.mapSeries(notAnsweredIssues, (issue) => addLabel(issue.url, NOT_ANSWERED_LABEL_TEXT));
};

updateLabels();
