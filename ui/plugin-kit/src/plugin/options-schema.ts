import fs from 'node:fs'
import path from 'node:path'
import * as TJS from 'typescript-json-schema'
import logger from './logger'

export function buildOptionsSchema({ dtsPath, outDir }: { dtsPath: string; outDir: string }) {
  // optionally pass argument to schema generator
  const settings: TJS.PartialArgs = {
    required: true,
    ignoreErrors: true,
  }

  // optionally pass ts compiler options
  const compilerOptions: TJS.CompilerOptions = {
    strictNullChecks: true,
  }

  // optionally pass a base path
  const program = TJS.getProgramFromFiles([dtsPath], compilerOptions)

  const generator = TJS.buildGenerator(program, settings)

  // extract from dtsPath file by regex of `ArtalkPlugin<XXX>`
  const dtsFileContent = fs.readFileSync(dtsPath, 'utf-8')
  const symbolNameMatches = /ArtalkPlugin<(.*?)>/.exec(dtsFileContent)
  if (!symbolNameMatches || symbolNameMatches.length < 2) {
    logger.info(`Skip plugin options schema generation.`)
    return
  }

  // Get symbols for different types from generator.
  const symbolName = symbolNameMatches[1]
  const schema = generator?.getSchemaForSymbol(symbolName)

  // Save schema
  const schemaPath = path.resolve(outDir, `artalk-plugin-options.schema.json`)
  fs.writeFileSync(schemaPath, JSON.stringify(schema, null, 2), 'utf-8')
}
