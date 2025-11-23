<?xml version="1.0" encoding="UTF-8"?>
<xsl:stylesheet version="2.0"
                xmlns:html="http://www.w3.org/1999/xhtml"
                xmlns:image="http://www.google.com/schemas/sitemap-image/1.1"
                xmlns:sitemap="http://www.sitemaps.org/schemas/sitemap/0.9"
                xmlns:xsl="http://www.w3.org/1999/XSL/Transform">
  <xsl:output method="html" version="1.0" encoding="UTF-8" indent="yes"/>
  <xsl:template match="/">
    <html xmlns="http://www.w3.org/1999/xhtml" lang="en">
      <head>
        <meta charset="UTF-8"/>
        <meta name="viewport" content="width=device-width, initial-scale=1.0"/>
        <title>XML Sitemap - Mockly</title>
        <meta name="description" content="Complete sitemap of all pages and resources available on Mockly"/>
        <style type="text/css">
          * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
          }
          
          body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            background: linear-gradient(to bottom, #0f172a, #1e293b, #0f172a);
            color: #e2e8f0;
            min-height: 100vh;
            padding: 2rem 1rem;
            line-height: 1.6;
          }
          
          .container {
            max-width: 1200px;
            margin: 0 auto;
          }
          
          header {
            text-align: center;
            margin-bottom: 3rem;
            padding: 2rem 0;
          }
          
          h1 {
            font-size: 3rem;
            font-weight: 700;
            color: #ffffff;
            margin-bottom: 0.5rem;
            background: linear-gradient(to right, #60a5fa, #3b82f6);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
          }
          
          .subtitle {
            color: #94a3b8;
            font-size: 1.125rem;
            margin-top: 0.5rem;
          }
          
          .intro {
            background: rgba(30, 41, 59, 0.5);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(148, 163, 184, 0.1);
            padding: 1.5rem;
            border-radius: 0.75rem;
            margin-bottom: 2rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
          }
          
          .intro p {
            color: #cbd5e1;
            margin: 0;
            font-size: 1rem;
          }
          
          .intro strong {
            color: #60a5fa;
            font-weight: 600;
          }
          
          .stats {
            display: flex;
            gap: 2rem;
            flex-wrap: wrap;
            margin-bottom: 2rem;
            padding: 1.5rem;
            background: rgba(30, 41, 59, 0.5);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(148, 163, 184, 0.1);
            border-radius: 0.75rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
          }
          
          .stat-item {
            display: flex;
            flex-direction: column;
            gap: 0.25rem;
          }
          
          .stat-label {
            color: #94a3b8;
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
          }
          
          .stat-value {
            color: #60a5fa;
            font-size: 1.5rem;
            font-weight: 700;
          }
          
          .table-wrapper {
            background: rgba(30, 41, 59, 0.5);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(148, 163, 184, 0.1);
            border-radius: 0.75rem;
            overflow: hidden;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
          }
          
          table {
            width: 100%;
            border-collapse: collapse;
            background: transparent;
          }
          
          th {
            background: rgba(15, 23, 42, 0.8);
            color: #ffffff;
            text-align: left;
            padding: 1rem;
            font-weight: 600;
            font-size: 0.875rem;
            text-transform: uppercase;
            letter-spacing: 0.05em;
            border-bottom: 2px solid rgba(59, 130, 246, 0.3);
          }
          
          th:first-child {
            padding-left: 1.5rem;
          }
          
          th:last-child {
            padding-right: 1.5rem;
          }
          
          td {
            padding: 1rem;
            border-bottom: 1px solid rgba(148, 163, 184, 0.1);
            color: #cbd5e1;
          }
          
          td:first-child {
            padding-left: 1.5rem;
          }
          
          td:last-child {
            padding-right: 1.5rem;
          }
          
          tr {
            transition: background-color 0.2s ease;
          }
          
          tr:hover {
            background: rgba(59, 130, 246, 0.1);
          }
          
          tr:last-child td {
            border-bottom: none;
          }
          
          a {
            color: #60a5fa;
            text-decoration: none;
            transition: color 0.2s ease;
            word-break: break-all;
          }
          
          a:hover {
            color: #93c5fd;
            text-decoration: underline;
          }
          
          .badge {
            display: inline-block;
            padding: 0.25rem 0.75rem;
            border-radius: 9999px;
            font-size: 0.75rem;
            font-weight: 600;
            text-transform: uppercase;
            letter-spacing: 0.05em;
          }
          
          .badge-daily {
            background: rgba(34, 197, 94, 0.2);
            color: #4ade80;
          }
          
          .badge-weekly {
            background: rgba(59, 130, 246, 0.2);
            color: #60a5fa;
          }
          
          .badge-monthly {
            background: rgba(168, 85, 247, 0.2);
            color: #a78bfa;
          }
          
          .badge-yearly {
            background: rgba(236, 72, 153, 0.2);
            color: #f472b6;
          }
          
          .badge-never {
            background: rgba(107, 114, 128, 0.2);
            color: #9ca3af;
          }
          
          .priority {
            font-weight: 600;
            color: #fbbf24;
          }
          
          .footer {
            margin-top: 3rem;
            padding: 2rem;
            text-align: center;
            background: rgba(30, 41, 59, 0.5);
            backdrop-filter: blur(10px);
            border: 1px solid rgba(148, 163, 184, 0.1);
            border-radius: 0.75rem;
            box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1);
          }
          
          .footer p {
            color: #94a3b8;
            margin-bottom: 0.5rem;
          }
          
          .footer a {
            color: #60a5fa;
            text-decoration: none;
            font-weight: 600;
          }
          
          .footer a:hover {
            text-decoration: underline;
          }
          
          @media (max-width: 768px) {
            body {
              padding: 1rem 0.5rem;
            }
            
            h1 {
              font-size: 2rem;
            }
            
            .stats {
              flex-direction: column;
              gap: 1rem;
            }
            
            table {
              font-size: 0.875rem;
            }
            
            th, td {
              padding: 0.75rem 0.5rem;
            }
            
            th:first-child, td:first-child {
              padding-left: 1rem;
            }
            
            th:last-child, td:last-child {
              padding-right: 1rem;
            }
          }
        </style>
      </head>
      <body>
        <div class="container">
          <header>
            <h1>XML Sitemap</h1>
            <p class="subtitle">Mockly - Schema-Driven Mock API</p>
          </header>
          
          <div class="intro">
            <p>
              This is an XML Sitemap generated by <strong>Mockly</strong>.
              It contains all available pages and resources on the website.
              Search engines use this file to discover and index content.
            </p>
          </div>
          
          <div class="stats">
            <div class="stat-item">
              <span class="stat-label">Total URLs</span>
              <span class="stat-value"><xsl:value-of select="count(sitemap:urlset/sitemap:url)"/></span>
            </div>
            <div class="stat-item">
              <span class="stat-label">Last Updated</span>
              <span class="stat-value">
                <xsl:value-of select="sitemap:urlset/sitemap:url[1]/sitemap:lastmod"/>
              </span>
            </div>
          </div>
          
          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>URL</th>
                  <th>Last Modified</th>
                  <th>Change Frequency</th>
                  <th>Priority</th>
                </tr>
              </thead>
              <tbody>
                <xsl:for-each select="sitemap:urlset/sitemap:url">
                  <xsl:sort select="sitemap:loc"/>
                  <tr>
                    <td>
                      <xsl:variable name="itemURL">
                        <xsl:value-of select="sitemap:loc"/>
                      </xsl:variable>
                      <a href="{$itemURL}">
                        <xsl:value-of select="sitemap:loc"/>
                      </a>
                    </td>
                    <td>
                      <xsl:value-of select="sitemap:lastmod"/>
                    </td>
                    <td>
                      <xsl:variable name="freq" select="sitemap:changefreq"/>
                      <span class="badge badge-{$freq}">
                        <xsl:value-of select="$freq"/>
                      </span>
                    </td>
                    <td>
                      <span class="priority">
                        <xsl:value-of select="sitemap:priority"/>
                      </span>
                    </td>
                  </tr>
                </xsl:for-each>
              </tbody>
            </table>
          </div>
          
          <div class="footer">
            <p>
              Generated by <strong>Mockly</strong> - Free Mock API Service
            </p>
            <p>
              <a href="/">Visit Homepage</a> | 
              <a href="/sitemap.html">View HTML Sitemap</a> | 
              <a href="/docs">Documentation</a>
            </p>
          </div>
        </div>
      </body>
    </html>
  </xsl:template>
</xsl:stylesheet>

