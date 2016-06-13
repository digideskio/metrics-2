package grammar

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/square/metrics/function"
)

var blahlbah = 3

var g = &grammar{
	rules: []*rule{
		{
			name: "Root",
			pos:  position{line: 8, col: 1, offset: 40},
			expr: &actionExpr{
				pos: position{line: 8, col: 9, offset: 48},
				run: (*parser).callonRoot1,
				expr: &seqExpr{
					pos: position{line: 8, col: 9, offset: 48},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 8, col: 9, offset: 48},
							label: "expr",
							expr: &choiceExpr{
								pos: position{line: 8, col: 15, offset: 54},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 8, col: 15, offset: 54},
										name: "SelectStatement",
									},
									&ruleRefExpr{
										pos:  position{line: 8, col: 33, offset: 72},
										name: "DescribeStatement",
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 8, col: 52, offset: 91},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 8, col: 54, offset: 93},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "SelectStatement",
			pos:  position{line: 12, col: 1, offset: 121},
			expr: &actionExpr{
				pos: position{line: 12, col: 20, offset: 140},
				run: (*parser).callonSelectStatement1,
				expr: &seqExpr{
					pos: position{line: 12, col: 20, offset: 140},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 12, col: 20, offset: 140},
							name: "_",
						},
						&zeroOrOneExpr{
							pos: position{line: 12, col: 22, offset: 142},
							expr: &seqExpr{
								pos: position{line: 12, col: 23, offset: 143},
								exprs: []interface{}{
									&litMatcher{
										pos:        position{line: 12, col: 23, offset: 143},
										val:        "select",
										ignoreCase: false,
									},
									&ruleRefExpr{
										pos:  position{line: 12, col: 32, offset: 152},
										name: "KEY",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 12, col: 38, offset: 158},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 43, offset: 163},
								name: "ExpressionList",
							},
						},
						&labeledExpr{
							pos:   position{line: 12, col: 58, offset: 178},
							label: "predicateClause",
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 74, offset: 194},
								name: "OptionalPredicateClause",
							},
						},
						&labeledExpr{
							pos:   position{line: 12, col: 98, offset: 218},
							label: "propertyClause",
							expr: &ruleRefExpr{
								pos:  position{line: 12, col: 113, offset: 233},
								name: "PropertyClause",
							},
						},
					},
				},
			},
		},
		{
			name: "DescribeStatement",
			pos:  position{line: 16, col: 1, offset: 332},
			expr: &actionExpr{
				pos: position{line: 17, col: 3, offset: 355},
				run: (*parser).callonDescribeStatement1,
				expr: &seqExpr{
					pos: position{line: 17, col: 3, offset: 355},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 17, col: 3, offset: 355},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 17, col: 5, offset: 357},
							val:        "describe",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 17, col: 16, offset: 368},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 18, col: 3, offset: 374},
							label: "statement",
							expr: &choiceExpr{
								pos: position{line: 18, col: 14, offset: 385},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 18, col: 14, offset: 385},
										name: "DescribeAllStatement",
									},
									&ruleRefExpr{
										pos:  position{line: 18, col: 37, offset: 408},
										name: "DescribeMetricsStatement",
									},
									&ruleRefExpr{
										pos:  position{line: 18, col: 64, offset: 435},
										name: "DescribeSingleStatement",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "DescribeAllStatement",
			pos:  position{line: 22, col: 1, offset: 489},
			expr: &actionExpr{
				pos: position{line: 23, col: 3, offset: 515},
				run: (*parser).callonDescribeAllStatement1,
				expr: &seqExpr{
					pos: position{line: 23, col: 3, offset: 515},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 23, col: 3, offset: 515},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 23, col: 5, offset: 517},
							val:        "all",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 23, col: 11, offset: 523},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 24, col: 3, offset: 529},
							label: "matchClause",
							expr: &ruleRefExpr{
								pos:  position{line: 24, col: 15, offset: 541},
								name: "OptionalMatchClause",
							},
						},
					},
				},
			},
		},
		{
			name: "OptionalMatchClause",
			pos:  position{line: 28, col: 1, offset: 629},
			expr: &choiceExpr{
				pos: position{line: 28, col: 24, offset: 652},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 28, col: 24, offset: 652},
						name: "MatchClause",
					},
					&ruleRefExpr{
						pos:  position{line: 28, col: 38, offset: 666},
						name: "EmptyMatchClause",
					},
				},
			},
		},
		{
			name: "EmptyMatchClause",
			pos:  position{line: 30, col: 1, offset: 684},
			expr: &actionExpr{
				pos: position{line: 30, col: 21, offset: 704},
				run: (*parser).callonEmptyMatchClause1,
				expr: &litMatcher{
					pos:        position{line: 30, col: 21, offset: 704},
					val:        "",
					ignoreCase: false,
				},
			},
		},
		{
			name: "MatchClause",
			pos:  position{line: 34, col: 1, offset: 748},
			expr: &actionExpr{
				pos: position{line: 35, col: 3, offset: 765},
				run: (*parser).callonMatchClause1,
				expr: &seqExpr{
					pos: position{line: 35, col: 3, offset: 765},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 35, col: 3, offset: 765},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 35, col: 5, offset: 767},
							val:        "match",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 35, col: 13, offset: 775},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 36, col: 3, offset: 781},
							label: "literal",
							expr: &ruleRefExpr{
								pos:  position{line: 36, col: 11, offset: 789},
								name: "LiteralString",
							},
						},
					},
				},
			},
		},
		{
			name: "DescribeMetrics",
			pos:  position{line: 40, col: 1, offset: 867},
			expr: &actionExpr{
				pos: position{line: 41, col: 3, offset: 888},
				run: (*parser).callonDescribeMetrics1,
				expr: &seqExpr{
					pos: position{line: 41, col: 3, offset: 888},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 41, col: 3, offset: 888},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 41, col: 5, offset: 890},
							val:        "metrics",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 41, col: 15, offset: 900},
							name: "KEY",
						},
						&ruleRefExpr{
							pos:  position{line: 42, col: 3, offset: 906},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 42, col: 5, offset: 908},
							val:        "where",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 42, col: 13, offset: 916},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 43, col: 3, offset: 922},
							label: "tagName",
							expr: &ruleRefExpr{
								pos:  position{line: 43, col: 11, offset: 930},
								name: "MetricName",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 44, col: 3, offset: 943},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 44, col: 5, offset: 945},
							val:        "=",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 45, col: 3, offset: 951},
							label: "tagValue",
							expr: &ruleRefExpr{
								pos:  position{line: 45, col: 12, offset: 960},
								name: "LiteralString",
							},
						},
					},
				},
			},
		},
		{
			name: "DescribeSingleStatement",
			pos:  position{line: 49, col: 1, offset: 1056},
			expr: &actionExpr{
				pos: position{line: 50, col: 3, offset: 1085},
				run: (*parser).callonDescribeSingleStatement1,
				expr: &seqExpr{
					pos: position{line: 50, col: 3, offset: 1085},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 50, col: 3, offset: 1085},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 50, col: 8, offset: 1090},
								name: "MetricName",
							},
						},
						&labeledExpr{
							pos:   position{line: 51, col: 3, offset: 1138},
							label: "predicateClause",
							expr: &ruleRefExpr{
								pos:  position{line: 51, col: 19, offset: 1154},
								name: "OptionalPredicateClause",
							},
						},
					},
				},
			},
		},
		{
			name: "UncheckedPropertyClause",
			pos:  position{line: 55, col: 1, offset: 1280},
			expr: &actionExpr{
				pos: position{line: 56, col: 3, offset: 1309},
				run: (*parser).callonUncheckedPropertyClause1,
				expr: &labeledExpr{
					pos:   position{line: 56, col: 3, offset: 1309},
					label: "propertyList",
					expr: &zeroOrMoreExpr{
						pos: position{line: 56, col: 16, offset: 1322},
						expr: &actionExpr{
							pos: position{line: 57, col: 7, offset: 1330},
							run: (*parser).callonUncheckedPropertyClause4,
							expr: &seqExpr{
								pos: position{line: 57, col: 7, offset: 1330},
								exprs: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 57, col: 7, offset: 1330},
										name: "_",
									},
									&labeledExpr{
										pos:   position{line: 57, col: 9, offset: 1332},
										label: "key",
										expr: &ruleRefExpr{
											pos:  position{line: 57, col: 13, offset: 1336},
											name: "PropertyKey",
										},
									},
									&ruleRefExpr{
										pos:  position{line: 58, col: 7, offset: 1354},
										name: "_",
									},
									&labeledExpr{
										pos:   position{line: 58, col: 9, offset: 1356},
										label: "value",
										expr: &ruleRefExpr{
											pos:  position{line: 58, col: 15, offset: 1362},
											name: "PropertyValue",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PropertyClause",
			pos:  position{line: 64, col: 1, offset: 1544},
			expr: &actionExpr{
				pos: position{line: 64, col: 19, offset: 1562},
				run: (*parser).callonPropertyClause1,
				expr: &labeledExpr{
					pos:   position{line: 64, col: 19, offset: 1562},
					label: "clause",
					expr: &ruleRefExpr{
						pos:  position{line: 64, col: 26, offset: 1569},
						name: "UncheckedPropertyClause",
					},
				},
			},
		},
		{
			name: "OptionalPredicateClause",
			pos:  position{line: 68, col: 1, offset: 1664},
			expr: &choiceExpr{
				pos: position{line: 68, col: 28, offset: 1691},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 68, col: 28, offset: 1691},
						name: "PredicateClause",
					},
					&actionExpr{
						pos: position{line: 68, col: 46, offset: 1709},
						run: (*parser).callonOptionalPredicateClause3,
						expr: &litMatcher{
							pos:        position{line: 68, col: 46, offset: 1709},
							val:        "",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExpressionList",
			pos:  position{line: 70, col: 1, offset: 1767},
			expr: &actionExpr{
				pos: position{line: 71, col: 3, offset: 1787},
				run: (*parser).callonExpressionList1,
				expr: &seqExpr{
					pos: position{line: 71, col: 3, offset: 1787},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 71, col: 3, offset: 1787},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 71, col: 9, offset: 1793},
								name: "Expression",
							},
						},
						&labeledExpr{
							pos:   position{line: 72, col: 3, offset: 1806},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 72, col: 9, offset: 1812},
								expr: &actionExpr{
									pos: position{line: 72, col: 11, offset: 1814},
									run: (*parser).callonExpressionList7,
									expr: &seqExpr{
										pos: position{line: 72, col: 11, offset: 1814},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 72, col: 11, offset: 1814},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 72, col: 13, offset: 1816},
												val:        ",",
												ignoreCase: false,
											},
											&labeledExpr{
												pos:   position{line: 72, col: 17, offset: 1820},
												label: "expression",
												expr: &ruleRefExpr{
													pos:  position{line: 72, col: 28, offset: 1831},
													name: "Expression",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 77, col: 1, offset: 1964},
			expr: &actionExpr{
				pos: position{line: 77, col: 15, offset: 1978},
				run: (*parser).callonExpression1,
				expr: &seqExpr{
					pos: position{line: 77, col: 15, offset: 1978},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 77, col: 15, offset: 1978},
							label: "sum",
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 19, offset: 1982},
								name: "ExpressionSum",
							},
						},
						&labeledExpr{
							pos:   position{line: 77, col: 33, offset: 1996},
							label: "pipe",
							expr: &ruleRefExpr{
								pos:  position{line: 77, col: 38, offset: 2001},
								name: "AddPipe",
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionSum",
			pos:  position{line: 79, col: 1, offset: 2070},
			expr: &actionExpr{
				pos: position{line: 80, col: 3, offset: 2089},
				run: (*parser).callonExpressionSum1,
				expr: &seqExpr{
					pos: position{line: 80, col: 3, offset: 2089},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 80, col: 3, offset: 2089},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 80, col: 9, offset: 2095},
								name: "ExpressionProduct",
							},
						},
						&labeledExpr{
							pos:   position{line: 81, col: 3, offset: 2115},
							label: "suffixes",
							expr: &zeroOrMoreExpr{
								pos: position{line: 81, col: 12, offset: 2124},
								expr: &actionExpr{
									pos: position{line: 82, col: 5, offset: 2130},
									run: (*parser).callonExpressionSum7,
									expr: &seqExpr{
										pos: position{line: 82, col: 5, offset: 2130},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 82, col: 5, offset: 2130},
												label: "pipes",
												expr: &zeroOrMoreExpr{
													pos: position{line: 82, col: 11, offset: 2136},
													expr: &ruleRefExpr{
														pos:  position{line: 82, col: 11, offset: 2136},
														name: "OnePipe",
													},
												},
											},
											&labeledExpr{
												pos:   position{line: 83, col: 5, offset: 2149},
												label: "op",
												expr: &actionExpr{
													pos: position{line: 83, col: 9, offset: 2153},
													run: (*parser).callonExpressionSum13,
													expr: &seqExpr{
														pos: position{line: 83, col: 9, offset: 2153},
														exprs: []interface{}{
															&ruleRefExpr{
																pos:  position{line: 83, col: 9, offset: 2153},
																name: "_",
															},
															&labeledExpr{
																pos:   position{line: 83, col: 11, offset: 2155},
																label: "op",
																expr: &choiceExpr{
																	pos: position{line: 83, col: 15, offset: 2159},
																	alternatives: []interface{}{
																		&litMatcher{
																			pos:        position{line: 83, col: 15, offset: 2159},
																			val:        "+",
																			ignoreCase: false,
																		},
																		&litMatcher{
																			pos:        position{line: 83, col: 21, offset: 2165},
																			val:        "-",
																			ignoreCase: false,
																		},
																	},
																},
															},
														},
													},
												},
											},
											&labeledExpr{
												pos:   position{line: 84, col: 5, offset: 2192},
												label: "right",
												expr: &ruleRefExpr{
													pos:  position{line: 84, col: 11, offset: 2198},
													name: "ExpressionProduct",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionProduct",
			pos:  position{line: 89, col: 1, offset: 2377},
			expr: &actionExpr{
				pos: position{line: 90, col: 3, offset: 2400},
				run: (*parser).callonExpressionProduct1,
				expr: &seqExpr{
					pos: position{line: 90, col: 3, offset: 2400},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 90, col: 3, offset: 2400},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 90, col: 9, offset: 2406},
								name: "ExpressionAtom",
							},
						},
						&labeledExpr{
							pos:   position{line: 91, col: 3, offset: 2423},
							label: "suffixes",
							expr: &zeroOrMoreExpr{
								pos: position{line: 91, col: 12, offset: 2432},
								expr: &actionExpr{
									pos: position{line: 92, col: 5, offset: 2438},
									run: (*parser).callonExpressionProduct7,
									expr: &seqExpr{
										pos: position{line: 92, col: 5, offset: 2438},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 92, col: 5, offset: 2438},
												label: "pipes",
												expr: &zeroOrMoreExpr{
													pos: position{line: 92, col: 11, offset: 2444},
													expr: &ruleRefExpr{
														pos:  position{line: 92, col: 11, offset: 2444},
														name: "OnePipe",
													},
												},
											},
											&labeledExpr{
												pos:   position{line: 93, col: 5, offset: 2457},
												label: "op",
												expr: &actionExpr{
													pos: position{line: 93, col: 9, offset: 2461},
													run: (*parser).callonExpressionProduct13,
													expr: &seqExpr{
														pos: position{line: 93, col: 9, offset: 2461},
														exprs: []interface{}{
															&ruleRefExpr{
																pos:  position{line: 93, col: 9, offset: 2461},
																name: "_",
															},
															&labeledExpr{
																pos:   position{line: 93, col: 11, offset: 2463},
																label: "op",
																expr: &choiceExpr{
																	pos: position{line: 93, col: 15, offset: 2467},
																	alternatives: []interface{}{
																		&litMatcher{
																			pos:        position{line: 93, col: 15, offset: 2467},
																			val:        "*",
																			ignoreCase: false,
																		},
																		&litMatcher{
																			pos:        position{line: 93, col: 21, offset: 2473},
																			val:        "/",
																			ignoreCase: false,
																		},
																	},
																},
															},
														},
													},
												},
											},
											&labeledExpr{
												pos:   position{line: 94, col: 5, offset: 2500},
												label: "right",
												expr: &ruleRefExpr{
													pos:  position{line: 94, col: 11, offset: 2506},
													name: "ExpressionAtom",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "OnePipe",
			pos:  position{line: 99, col: 1, offset: 2676},
			expr: &actionExpr{
				pos: position{line: 100, col: 3, offset: 2689},
				run: (*parser).callonOnePipe1,
				expr: &seqExpr{
					pos: position{line: 100, col: 3, offset: 2689},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 100, col: 3, offset: 2689},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 100, col: 5, offset: 2691},
							val:        "|",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 101, col: 3, offset: 2697},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 101, col: 5, offset: 2699},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 101, col: 10, offset: 2704},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 102, col: 3, offset: 2737},
							label: "arguments",
							expr: &choiceExpr{
								pos: position{line: 102, col: 14, offset: 2748},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 102, col: 14, offset: 2748},
										name: "CallArguments",
									},
									&actionExpr{
										pos: position{line: 102, col: 30, offset: 2764},
										run: (*parser).callonOnePipe11,
										expr: &litMatcher{
											pos:        position{line: 102, col: 30, offset: 2764},
											val:        "",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CallArguments",
			pos:  position{line: 105, col: 1, offset: 2852},
			expr: &actionExpr{
				pos: position{line: 106, col: 3, offset: 2871},
				run: (*parser).callonCallArguments1,
				expr: &seqExpr{
					pos: position{line: 106, col: 3, offset: 2871},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 106, col: 3, offset: 2871},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 106, col: 5, offset: 2873},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 107, col: 3, offset: 2879},
							label: "arguments",
							expr: &choiceExpr{
								pos: position{line: 107, col: 14, offset: 2890},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 107, col: 15, offset: 2891},
										exprs: []interface{}{
											&labeledExpr{
												pos:   position{line: 107, col: 15, offset: 2891},
												label: "first",
												expr: &ruleRefExpr{
													pos:  position{line: 107, col: 21, offset: 2897},
													name: "Expression",
												},
											},
											&labeledExpr{
												pos:   position{line: 107, col: 32, offset: 2908},
												label: "rest",
												expr: &zeroOrMoreExpr{
													pos: position{line: 107, col: 37, offset: 2913},
													expr: &seqExpr{
														pos: position{line: 107, col: 38, offset: 2914},
														exprs: []interface{}{
															&ruleRefExpr{
																pos:  position{line: 107, col: 38, offset: 2914},
																name: "_",
															},
															&litMatcher{
																pos:        position{line: 107, col: 40, offset: 2916},
																val:        ",",
																ignoreCase: false,
															},
															&labeledExpr{
																pos:   position{line: 107, col: 44, offset: 2920},
																label: "right",
																expr: &ruleRefExpr{
																	pos:  position{line: 107, col: 50, offset: 2926},
																	name: "Expression",
																},
															},
														},
													},
												},
											},
										},
									},
									&actionExpr{
										pos: position{line: 107, col: 66, offset: 2942},
										run: (*parser).callonCallArguments17,
										expr: &litMatcher{
											pos:        position{line: 107, col: 66, offset: 2942},
											val:        "",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 109, col: 3, offset: 2988},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 109, col: 5, offset: 2990},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExpressionAtom",
			pos:  position{line: 112, col: 1, offset: 3023},
			expr: &actionExpr{
				pos: position{line: 113, col: 3, offset: 3043},
				run: (*parser).callonExpressionAtom1,
				expr: &seqExpr{
					pos: position{line: 113, col: 3, offset: 3043},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 113, col: 3, offset: 3043},
							label: "core",
							expr: &ruleRefExpr{
								pos:  position{line: 113, col: 8, offset: 3048},
								name: "ExpressionRaw",
							},
						},
						&labeledExpr{
							pos:   position{line: 114, col: 3, offset: 3064},
							label: "annotation",
							expr: &ruleRefExpr{
								pos:  position{line: 114, col: 14, offset: 3075},
								name: "ExpressionAnnotation",
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionRaw",
			pos:  position{line: 121, col: 1, offset: 3231},
			expr: &choiceExpr{
				pos: position{line: 121, col: 18, offset: 3248},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 121, col: 18, offset: 3248},
						name: "ExpressionFunction",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 39, offset: 3269},
						name: "ExpressionMetric",
					},
					&actionExpr{
						pos: position{line: 121, col: 58, offset: 3288},
						run: (*parser).callonExpressionRaw4,
						expr: &seqExpr{
							pos: position{line: 121, col: 58, offset: 3288},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 121, col: 58, offset: 3288},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 121, col: 60, offset: 3290},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 121, col: 64, offset: 3294},
									label: "item",
									expr: &ruleRefExpr{
										pos:  position{line: 121, col: 69, offset: 3299},
										name: "Expression",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 121, col: 80, offset: 3310},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 121, col: 82, offset: 3312},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 107, offset: 3337},
						name: "Duration",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 118, offset: 3348},
						name: "Scalar",
					},
					&ruleRefExpr{
						pos:  position{line: 121, col: 127, offset: 3357},
						name: "String",
					},
				},
			},
		},
		{
			name: "ExpressionAnnotationRequired",
			pos:  position{line: 123, col: 1, offset: 3365},
			expr: &actionExpr{
				pos: position{line: 123, col: 33, offset: 3397},
				run: (*parser).callonExpressionAnnotationRequired1,
				expr: &seqExpr{
					pos: position{line: 123, col: 33, offset: 3397},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 123, col: 33, offset: 3397},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 123, col: 35, offset: 3399},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 123, col: 39, offset: 3403},
							label: "contents",
							expr: &charClassMatcher{
								pos:        position{line: 123, col: 48, offset: 3412},
								val:        "[^}]",
								chars:      []rune{'}'},
								ignoreCase: false,
								inverted:   true,
							},
						},
						&litMatcher{
							pos:        position{line: 123, col: 53, offset: 3417},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExpressionAnnotation",
			pos:  position{line: 124, col: 1, offset: 3446},
			expr: &choiceExpr{
				pos: position{line: 124, col: 25, offset: 3470},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 124, col: 25, offset: 3470},
						name: "ExpressionAnnotationRequires",
					},
					&actionExpr{
						pos: position{line: 124, col: 56, offset: 3501},
						run: (*parser).callonExpressionAnnotation3,
						expr: &litMatcher{
							pos:        position{line: 124, col: 56, offset: 3501},
							val:        "",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "OptionalGroupBy",
			pos:  position{line: 126, col: 1, offset: 3524},
			expr: &choiceExpr{
				pos: position{line: 126, col: 20, offset: 3543},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 126, col: 20, offset: 3543},
						name: "GroupByClause",
					},
					&ruleRefExpr{
						pos:  position{line: 126, col: 36, offset: 3559},
						name: "CollapseByClause",
					},
					&actionExpr{
						pos: position{line: 126, col: 55, offset: 3578},
						run: (*parser).callonOptionalGroupBy4,
						expr: &litMatcher{
							pos:        position{line: 126, col: 55, offset: 3578},
							val:        "",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ExpressionFunction",
			pos:  position{line: 128, col: 1, offset: 3645},
			expr: &actionExpr{
				pos: position{line: 129, col: 3, offset: 3669},
				run: (*parser).callonExpressionFunction1,
				expr: &seqExpr{
					pos: position{line: 129, col: 3, offset: 3669},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 129, col: 3, offset: 3669},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 129, col: 8, offset: 3674},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 130, col: 3, offset: 3687},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 13, offset: 3697},
								name: "CallArguments",
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionMetric",
			pos:  position{line: 133, col: 1, offset: 3785},
			expr: &actionExpr{
				pos: position{line: 134, col: 3, offset: 3807},
				run: (*parser).callonExpressionMetric1,
				expr: &seqExpr{
					pos: position{line: 134, col: 3, offset: 3807},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 134, col: 3, offset: 3807},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 134, col: 8, offset: 3812},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 135, col: 3, offset: 3825},
							label: "predicate",
							expr: &choiceExpr{
								pos: position{line: 135, col: 15, offset: 3837},
								alternatives: []interface{}{
									&actionExpr{
										pos: position{line: 135, col: 15, offset: 3837},
										run: (*parser).callonExpressionMetric7,
										expr: &seqExpr{
											pos: position{line: 135, col: 15, offset: 3837},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 135, col: 15, offset: 3837},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 135, col: 17, offset: 3839},
													val:        "[",
													ignoreCase: false,
												},
												&labeledExpr{
													pos:   position{line: 135, col: 21, offset: 3843},
													label: "predicate",
													expr: &ruleRefExpr{
														pos:  position{line: 135, col: 31, offset: 3853},
														name: "Predicate",
													},
												},
												&ruleRefExpr{
													pos:  position{line: 135, col: 41, offset: 3863},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 135, col: 43, offset: 3865},
													val:        "]",
													ignoreCase: false,
												},
											},
										},
									},
									&actionExpr{
										pos: position{line: 135, col: 75, offset: 3897},
										run: (*parser).callonExpressionMetric15,
										expr: &litMatcher{
											pos:        position{line: 135, col: 75, offset: 3897},
											val:        "",
											ignoreCase: false,
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "GroupByClause",
			pos:  position{line: 138, col: 1, offset: 4014},
			expr: &actionExpr{
				pos: position{line: 139, col: 3, offset: 4033},
				run: (*parser).callonGroupByClause1,
				expr: &seqExpr{
					pos: position{line: 139, col: 3, offset: 4033},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 139, col: 3, offset: 4033},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 139, col: 5, offset: 4035},
							val:        "group",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 139, col: 13, offset: 4043},
							name: "KEY",
						},
						&ruleRefExpr{
							pos:  position{line: 140, col: 3, offset: 4049},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 140, col: 5, offset: 4051},
							val:        "by",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 140, col: 10, offset: 4056},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 141, col: 3, offset: 4062},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 9, offset: 4068},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 142, col: 3, offset: 4142},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 142, col: 8, offset: 4147},
								expr: &seqExpr{
									pos: position{line: 142, col: 9, offset: 4148},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 142, col: 9, offset: 4148},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 142, col: 11, offset: 4150},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 142, col: 15, offset: 4154},
											name: "Identifier",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "CollapseByClause",
			pos:  position{line: 145, col: 1, offset: 4227},
			expr: &actionExpr{
				pos: position{line: 146, col: 3, offset: 4249},
				run: (*parser).callonCollapseByClause1,
				expr: &seqExpr{
					pos: position{line: 146, col: 3, offset: 4249},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 146, col: 3, offset: 4249},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 146, col: 5, offset: 4251},
							val:        "collapse",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 146, col: 16, offset: 4262},
							name: "KEY",
						},
						&ruleRefExpr{
							pos:  position{line: 147, col: 3, offset: 4268},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 147, col: 5, offset: 4270},
							val:        "by",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 147, col: 10, offset: 4275},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 148, col: 3, offset: 4281},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 148, col: 9, offset: 4287},
								name: "Identifier",
							},
						},
						&labeledExpr{
							pos:   position{line: 149, col: 3, offset: 4361},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 149, col: 8, offset: 4366},
								expr: &seqExpr{
									pos: position{line: 149, col: 9, offset: 4367},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 149, col: 9, offset: 4367},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 149, col: 11, offset: 4369},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 149, col: 15, offset: 4373},
											name: "Identifier",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PredicateClause",
			pos:  position{line: 152, col: 1, offset: 4449},
			expr: &actionExpr{
				pos: position{line: 153, col: 3, offset: 4470},
				run: (*parser).callonPredicateClause1,
				expr: &seqExpr{
					pos: position{line: 153, col: 3, offset: 4470},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 153, col: 3, offset: 4470},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 153, col: 5, offset: 4472},
							val:        "where",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 153, col: 13, offset: 4480},
							name: "KEY",
						},
						&labeledExpr{
							pos:   position{line: 154, col: 3, offset: 4486},
							label: "predicate",
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 13, offset: 4496},
								name: "Predicate",
							},
						},
					},
				},
			},
		},
		{
			name: "Predicate",
			pos:  position{line: 157, col: 1, offset: 4535},
			expr: &ruleRefExpr{
				pos:  position{line: 157, col: 14, offset: 4548},
				name: "PredicateDisjunction",
			},
		},
		{
			name: "PredicateDisjunction",
			pos:  position{line: 159, col: 1, offset: 4570},
			expr: &choiceExpr{
				pos: position{line: 160, col: 3, offset: 4596},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 160, col: 3, offset: 4596},
						run: (*parser).callonPredicateDisjunction2,
						expr: &seqExpr{
							pos: position{line: 160, col: 3, offset: 4596},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 160, col: 3, offset: 4596},
									label: "left",
									expr: &ruleRefExpr{
										pos:  position{line: 160, col: 8, offset: 4601},
										name: "PredicateConjunction",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 160, col: 29, offset: 4622},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 160, col: 31, offset: 4624},
									val:        "or",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 160, col: 36, offset: 4629},
									name: "KEY",
								},
								&labeledExpr{
									pos:   position{line: 160, col: 40, offset: 4633},
									label: "right",
									expr: &ruleRefExpr{
										pos:  position{line: 160, col: 46, offset: 4639},
										name: "PredicateDisjunction",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 162, col: 3, offset: 4731},
						name: "PredicateConjunction",
					},
				},
			},
		},
		{
			name: "PredicateConjunction",
			pos:  position{line: 164, col: 1, offset: 4753},
			expr: &choiceExpr{
				pos: position{line: 165, col: 3, offset: 4779},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 165, col: 3, offset: 4779},
						run: (*parser).callonPredicateConjunction2,
						expr: &seqExpr{
							pos: position{line: 165, col: 3, offset: 4779},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 165, col: 3, offset: 4779},
									label: "left",
									expr: &ruleRefExpr{
										pos:  position{line: 165, col: 8, offset: 4784},
										name: "PredicateAtom",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 165, col: 22, offset: 4798},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 165, col: 24, offset: 4800},
									val:        "and",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 165, col: 30, offset: 4806},
									name: "KEY",
								},
								&labeledExpr{
									pos:   position{line: 165, col: 34, offset: 4810},
									label: "right",
									expr: &ruleRefExpr{
										pos:  position{line: 165, col: 40, offset: 4816},
										name: "PredicateConjunction",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 167, col: 3, offset: 4910},
						name: "PredicateAtom",
					},
				},
			},
		},
		{
			name: "PredicateAtom",
			pos:  position{line: 169, col: 1, offset: 4925},
			expr: &choiceExpr{
				pos: position{line: 170, col: 3, offset: 4944},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 170, col: 3, offset: 4944},
						run: (*parser).callonPredicateAtom2,
						expr: &seqExpr{
							pos: position{line: 170, col: 3, offset: 4944},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 170, col: 3, offset: 4944},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 170, col: 5, offset: 4946},
									val:        "not",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 170, col: 11, offset: 4952},
									label: "atom",
									expr: &ruleRefExpr{
										pos:  position{line: 170, col: 16, offset: 4957},
										name: "PredicateAtom",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 172, col: 3, offset: 5037},
						run: (*parser).callonPredicateAtom8,
						expr: &seqExpr{
							pos: position{line: 172, col: 3, offset: 5037},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 172, col: 3, offset: 5037},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 172, col: 5, offset: 5039},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 172, col: 9, offset: 5043},
									label: "predicate",
									expr: &ruleRefExpr{
										pos:  position{line: 172, col: 19, offset: 5053},
										name: "Predicate",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 172, col: 29, offset: 5063},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 172, col: 31, offset: 5065},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 174, col: 3, offset: 5101},
						name: "TagMatcher",
					},
				},
			},
		},
		{
			name: "TagMatcher",
			pos:  position{line: 176, col: 1, offset: 5113},
			expr: &choiceExpr{
				pos: position{line: 177, col: 3, offset: 5129},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 177, col: 3, offset: 5129},
						run: (*parser).callonTagMatcher2,
						expr: &seqExpr{
							pos: position{line: 177, col: 3, offset: 5129},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 177, col: 3, offset: 5129},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 177, col: 8, offset: 5134},
										name: "Identifer",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 177, col: 18, offset: 5144},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 177, col: 20, offset: 5146},
									val:        "=",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 177, col: 24, offset: 5150},
									label: "literal",
									expr: &ruleRefExpr{
										pos:  position{line: 177, col: 32, offset: 5158},
										name: "LiteralString",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 179, col: 3, offset: 5290},
						run: (*parser).callonTagMatcher10,
						expr: &seqExpr{
							pos: position{line: 179, col: 3, offset: 5290},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 179, col: 3, offset: 5290},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 179, col: 8, offset: 5295},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 179, col: 19, offset: 5306},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 179, col: 21, offset: 5308},
									val:        "!=",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 179, col: 26, offset: 5313},
									label: "literal",
									expr: &ruleRefExpr{
										pos:  position{line: 179, col: 34, offset: 5321},
										name: "LiteralString",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 181, col: 3, offset: 5415},
						run: (*parser).callonTagMatcher18,
						expr: &seqExpr{
							pos: position{line: 181, col: 3, offset: 5415},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 181, col: 3, offset: 5415},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 181, col: 8, offset: 5420},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 181, col: 19, offset: 5431},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 181, col: 21, offset: 5433},
									val:        "match",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 181, col: 29, offset: 5441},
									name: "KEY",
								},
								&labeledExpr{
									pos:   position{line: 181, col: 33, offset: 5445},
									label: "literal",
									expr: &ruleRefExpr{
										pos:  position{line: 181, col: 41, offset: 5453},
										name: "LiteralString",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 183, col: 3, offset: 5547},
						run: (*parser).callonTagMatcher27,
						expr: &seqExpr{
							pos: position{line: 183, col: 3, offset: 5547},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 183, col: 3, offset: 5547},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 183, col: 8, offset: 5552},
										name: "Identifier",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 183, col: 19, offset: 5563},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 183, col: 21, offset: 5565},
									val:        "in",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 183, col: 26, offset: 5570},
									name: "KEY",
								},
								&labeledExpr{
									pos:   position{line: 183, col: 30, offset: 5574},
									label: "list",
									expr: &ruleRefExpr{
										pos:  position{line: 183, col: 35, offset: 5579},
										name: "LiteralList",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "LiteralString",
			pos:  position{line: 185, col: 1, offset: 5661},
			expr: &actionExpr{
				pos: position{line: 185, col: 18, offset: 5678},
				run: (*parser).callonLiteralString1,
				expr: &seqExpr{
					pos: position{line: 185, col: 18, offset: 5678},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 185, col: 18, offset: 5678},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 185, col: 20, offset: 5680},
							val:        "\"",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 185, col: 25, offset: 5685},
							label: "contents",
							expr: &zeroOrMoreExpr{
								pos: position{line: 185, col: 34, offset: 5694},
								expr: &choiceExpr{
									pos: position{line: 185, col: 35, offset: 5695},
									alternatives: []interface{}{
										&charClassMatcher{
											pos:        position{line: 185, col: 35, offset: 5695},
											val:        "[^\"]",
											chars:      []rune{'"'},
											ignoreCase: false,
											inverted:   true,
										},
										&litMatcher{
											pos:        position{line: 185, col: 42, offset: 5702},
											val:        "\\\"",
											ignoreCase: false,
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 185, col: 51, offset: 5711},
							val:        "\"",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "LiteralList",
			pos:  position{line: 187, col: 1, offset: 5770},
			expr: &actionExpr{
				pos: position{line: 188, col: 3, offset: 5787},
				run: (*parser).callonLiteralList1,
				expr: &seqExpr{
					pos: position{line: 188, col: 3, offset: 5787},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 188, col: 3, offset: 5787},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 188, col: 5, offset: 5789},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 189, col: 3, offset: 5795},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 189, col: 9, offset: 5801},
								name: "LiteralString",
							},
						},
						&labeledExpr{
							pos:   position{line: 190, col: 3, offset: 5817},
							label: "rest",
							expr: &actionExpr{
								pos: position{line: 190, col: 9, offset: 5823},
								run: (*parser).callonLiteralList8,
								expr: &seqExpr{
									pos: position{line: 190, col: 9, offset: 5823},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 190, col: 9, offset: 5823},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 190, col: 11, offset: 5825},
											val:        ",",
											ignoreCase: false,
										},
										&labeledExpr{
											pos:   position{line: 190, col: 15, offset: 5829},
											label: "literal",
											expr: &ruleRefExpr{
												pos:  position{line: 190, col: 23, offset: 5837},
												name: "LiteralString",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 191, col: 3, offset: 5878},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 191, col: 5, offset: 5880},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 194, col: 1, offset: 5957},
			expr: &choiceExpr{
				pos: position{line: 195, col: 3, offset: 5973},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 195, col: 3, offset: 5973},
						run: (*parser).callonIdentifier2,
						expr: &seqExpr{
							pos: position{line: 195, col: 3, offset: 5973},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 195, col: 3, offset: 5973},
									val:        "`",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 195, col: 7, offset: 5977},
									label: "contents",
									expr: &zeroOrMoreExpr{
										pos: position{line: 195, col: 16, offset: 5986},
										expr: &ruleRefExpr{
											pos:  position{line: 195, col: 16, offset: 5986},
											name: "CHAR",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 195, col: 22, offset: 5992},
									val:        "`",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 197, col: 3, offset: 6026},
						run: (*parser).callonIdentifier9,
						expr: &seqExpr{
							pos: position{line: 197, col: 3, offset: 6026},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 197, col: 3, offset: 6026},
									name: "_",
								},
								&notExpr{
									pos: position{line: 197, col: 5, offset: 6028},
									expr: &seqExpr{
										pos: position{line: 197, col: 7, offset: 6030},
										exprs: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 197, col: 7, offset: 6030},
												name: "KEYWORD",
											},
											&ruleRefExpr{
												pos:  position{line: 197, col: 15, offset: 6038},
												name: "KEY",
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 197, col: 20, offset: 6043},
									label: "first",
									expr: &ruleRefExpr{
										pos:  position{line: 197, col: 26, offset: 6049},
										name: "Segment",
									},
								},
								&labeledExpr{
									pos:   position{line: 197, col: 34, offset: 6057},
									label: "rest",
									expr: &zeroOrMoreExpr{
										pos: position{line: 197, col: 39, offset: 6062},
										expr: &actionExpr{
											pos: position{line: 197, col: 40, offset: 6063},
											run: (*parser).callonIdentifier20,
											expr: &seqExpr{
												pos: position{line: 197, col: 40, offset: 6063},
												exprs: []interface{}{
													&litMatcher{
														pos:        position{line: 197, col: 40, offset: 6063},
														val:        ".",
														ignoreCase: false,
													},
													&labeledExpr{
														pos:   position{line: 197, col: 44, offset: 6067},
														label: "segment",
														expr: &ruleRefExpr{
															pos:  position{line: 197, col: 52, offset: 6075},
															name: "Segment",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Segment",
			pos:  position{line: 199, col: 1, offset: 6178},
			expr: &seqExpr{
				pos: position{line: 199, col: 12, offset: 6189},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 199, col: 12, offset: 6189},
						name: "IDStart",
					},
					&zeroOrMoreExpr{
						pos: position{line: 199, col: 20, offset: 6197},
						expr: &ruleRefExpr{
							pos:  position{line: 199, col: 20, offset: 6197},
							name: "IDContinue",
						},
					},
				},
			},
		},
		{
			name: "IDStart",
			pos:  position{line: 201, col: 1, offset: 6210},
			expr: &charClassMatcher{
				pos:        position{line: 201, col: 12, offset: 6221},
				val:        "[a-zA-Z_]",
				chars:      []rune{'_'},
				ranges:     []rune{'a', 'z', 'A', 'Z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "IDContinue",
			pos:  position{line: 202, col: 1, offset: 6231},
			expr: &choiceExpr{
				pos: position{line: 202, col: 15, offset: 6245},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 202, col: 15, offset: 6245},
						name: "IDStart",
					},
					&charClassMatcher{
						pos:        position{line: 202, col: 25, offset: 6255},
						val:        "[0-9]",
						ranges:     []rune{'0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "Timestamp",
			pos:  position{line: 204, col: 1, offset: 6262},
			expr: &choiceExpr{
				pos: position{line: 204, col: 14, offset: 6275},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 204, col: 14, offset: 6275},
						run: (*parser).callonTimestamp2,
						expr: &seqExpr{
							pos: position{line: 204, col: 14, offset: 6275},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 204, col: 14, offset: 6275},
									label: "number",
									expr: &ruleRefExpr{
										pos:  position{line: 204, col: 21, offset: 6282},
										name: "NumberInteger",
									},
								},
								&labeledExpr{
									pos:   position{line: 204, col: 35, offset: 6296},
									label: "suffix",
									expr: &zeroOrMoreExpr{
										pos: position{line: 204, col: 42, offset: 6303},
										expr: &charClassMatcher{
											pos:        position{line: 204, col: 42, offset: 6303},
											val:        "[a-z]",
											ranges:     []rune{'a', 'z'},
											ignoreCase: false,
											inverted:   false,
										},
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 204, col: 83, offset: 6344},
						name: "LiteralString",
					},
					&actionExpr{
						pos: position{line: 204, col: 99, offset: 6360},
						run: (*parser).callonTimestamp10,
						expr: &seqExpr{
							pos: position{line: 204, col: 99, offset: 6360},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 204, col: 99, offset: 6360},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 204, col: 101, offset: 6362},
									val:        "now",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 204, col: 107, offset: 6368},
									name: "KEY",
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 206, col: 1, offset: 6395},
			expr: &zeroOrMoreExpr{
				pos: position{line: 206, col: 19, offset: 6413},
				expr: &charClassMatcher{
					pos:        position{line: 206, col: 19, offset: 6413},
					val:        "[ \\n\\t\\r]",
					chars:      []rune{' ', '\n', '\t', '\r'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 208, col: 1, offset: 6425},
			expr: &notExpr{
				pos: position{line: 208, col: 8, offset: 6432},
				expr: &anyMatcher{
					line: 208, col: 9, offset: 6433,
				},
			},
		},
		{
			name: "PropertyKey",
			pos:  position{line: 210, col: 1, offset: 6436},
			expr: &choiceExpr{
				pos: position{line: 210, col: 16, offset: 6451},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 210, col: 16, offset: 6451},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 210, col: 16, offset: 6451},
								val:        "from",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 210, col: 23, offset: 6458},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 210, col: 29, offset: 6464},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 210, col: 29, offset: 6464},
								val:        "to",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 210, col: 34, offset: 6469},
								name: "KEY",
							},
						},
					},
					&actionExpr{
						pos: position{line: 210, col: 40, offset: 6475},
						run: (*parser).callonPropertyKey8,
						expr: &seqExpr{
							pos: position{line: 210, col: 40, offset: 6475},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 210, col: 40, offset: 6475},
									val:        "sample",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 210, col: 49, offset: 6484},
									name: "KEY",
								},
								&ruleRefExpr{
									pos:  position{line: 210, col: 53, offset: 6488},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 210, col: 55, offset: 6490},
									val:        "by",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 210, col: 60, offset: 6495},
									name: "KEY",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "PropertyValue",
			pos:  position{line: 212, col: 1, offset: 6525},
			expr: &ruleRefExpr{
				pos:  position{line: 212, col: 18, offset: 6542},
				name: "Timestamp",
			},
		},
		{
			name: "KEYWORD",
			pos:  position{line: 214, col: 1, offset: 6553},
			expr: &choiceExpr{
				pos: position{line: 215, col: 3, offset: 6566},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 215, col: 3, offset: 6566},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 215, col: 3, offset: 6566},
								val:        "all",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 215, col: 9, offset: 6572},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 216, col: 3, offset: 6580},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 216, col: 3, offset: 6580},
								val:        "and",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 216, col: 9, offset: 6586},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 217, col: 3, offset: 6594},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 217, col: 3, offset: 6594},
								val:        "as",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 217, col: 8, offset: 6599},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 218, col: 3, offset: 6607},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 218, col: 3, offset: 6607},
								val:        "by",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 218, col: 8, offset: 6612},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 219, col: 3, offset: 6620},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 219, col: 3, offset: 6620},
								val:        "describe",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 219, col: 14, offset: 6631},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 220, col: 3, offset: 6639},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 220, col: 3, offset: 6639},
								val:        "group",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 220, col: 11, offset: 6647},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 221, col: 3, offset: 6655},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 221, col: 3, offset: 6655},
								val:        "collapse",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 221, col: 14, offset: 6666},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 222, col: 3, offset: 6674},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 222, col: 3, offset: 6674},
								val:        "in",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 222, col: 8, offset: 6679},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 223, col: 3, offset: 6687},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 223, col: 3, offset: 6687},
								val:        "match",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 223, col: 11, offset: 6695},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 224, col: 3, offset: 6703},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 224, col: 3, offset: 6703},
								val:        "not",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 224, col: 9, offset: 6709},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 225, col: 3, offset: 6717},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 225, col: 3, offset: 6717},
								val:        "or",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 225, col: 8, offset: 6722},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 226, col: 3, offset: 6730},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 226, col: 3, offset: 6730},
								val:        "select",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 226, col: 12, offset: 6739},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 227, col: 3, offset: 6747},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 227, col: 3, offset: 6747},
								val:        "where",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 227, col: 11, offset: 6755},
								name: "KEY",
							},
						},
					},
					&seqExpr{
						pos: position{line: 228, col: 3, offset: 6763},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 228, col: 3, offset: 6763},
								val:        "metrics",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 228, col: 13, offset: 6773},
								name: "KEY",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 229, col: 3, offset: 6781},
						name: "PROPERTY_KEY",
					},
				},
			},
		},
		{
			name: "Number",
			pos:  position{line: 231, col: 1, offset: 6795},
			expr: &seqExpr{
				pos: position{line: 231, col: 11, offset: 6805},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 231, col: 11, offset: 6805},
						name: "NumberInteger",
					},
					&zeroOrOneExpr{
						pos: position{line: 231, col: 25, offset: 6819},
						expr: &ruleRefExpr{
							pos:  position{line: 231, col: 25, offset: 6819},
							name: "NumberFraction",
						},
					},
					&zeroOrOneExpr{
						pos: position{line: 231, col: 41, offset: 6835},
						expr: &ruleRefExpr{
							pos:  position{line: 231, col: 41, offset: 6835},
							name: "NumberExponent",
						},
					},
				},
			},
		},
		{
			name: "NumberInteger",
			pos:  position{line: 232, col: 1, offset: 6851},
			expr: &seqExpr{
				pos: position{line: 232, col: 18, offset: 6868},
				exprs: []interface{}{
					&zeroOrOneExpr{
						pos: position{line: 232, col: 18, offset: 6868},
						expr: &litMatcher{
							pos:        position{line: 232, col: 18, offset: 6868},
							val:        "-",
							ignoreCase: false,
						},
					},
					&ruleRefExpr{
						pos:  position{line: 232, col: 23, offset: 6873},
						name: "NumberNatural",
					},
				},
			},
		},
		{
			name: "NumberNatural",
			pos:  position{line: 233, col: 1, offset: 6887},
			expr: &choiceExpr{
				pos: position{line: 233, col: 18, offset: 6904},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 233, col: 18, offset: 6904},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 233, col: 24, offset: 6910},
						exprs: []interface{}{
							&charClassMatcher{
								pos:        position{line: 233, col: 24, offset: 6910},
								val:        "[1-9]",
								ranges:     []rune{'1', '9'},
								ignoreCase: false,
								inverted:   false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 233, col: 30, offset: 6916},
								expr: &charClassMatcher{
									pos:        position{line: 233, col: 30, offset: 6916},
									val:        "[0-9]",
									ranges:     []rune{'0', '9'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NumberFraction",
			pos:  position{line: 234, col: 1, offset: 6923},
			expr: &seqExpr{
				pos: position{line: 234, col: 19, offset: 6941},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 234, col: 19, offset: 6941},
						val:        ".",
						ignoreCase: false,
					},
					&oneOrMoreExpr{
						pos: position{line: 234, col: 23, offset: 6945},
						expr: &charClassMatcher{
							pos:        position{line: 234, col: 23, offset: 6945},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name: "NumberExponent",
			pos:  position{line: 235, col: 1, offset: 6952},
			expr: &seqExpr{
				pos: position{line: 235, col: 19, offset: 6970},
				exprs: []interface{}{
					&charClassMatcher{
						pos:        position{line: 235, col: 19, offset: 6970},
						val:        "[eE]",
						chars:      []rune{'e', 'E'},
						ignoreCase: false,
						inverted:   false,
					},
					&zeroOrOneExpr{
						pos: position{line: 235, col: 24, offset: 6975},
						expr: &choiceExpr{
							pos: position{line: 235, col: 25, offset: 6976},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 235, col: 25, offset: 6976},
									val:        "+",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 235, col: 31, offset: 6982},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 235, col: 37, offset: 6988},
						expr: &charClassMatcher{
							pos:        position{line: 235, col: 37, offset: 6988},
							val:        "[0-9]",
							ranges:     []rune{'0', '9'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name: "Duration",
			pos:  position{line: 236, col: 1, offset: 6995},
			expr: &seqExpr{
				pos: position{line: 236, col: 13, offset: 7007},
				exprs: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 236, col: 13, offset: 7007},
						name: "NumberInteger",
					},
					&oneOrMoreExpr{
						pos: position{line: 236, col: 27, offset: 7021},
						expr: &charClassMatcher{
							pos:        position{line: 236, col: 27, offset: 7021},
							val:        "[a-z]",
							ranges:     []rune{'a', 'z'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name: "KEY",
			pos:  position{line: 238, col: 1, offset: 7029},
			expr: &notExpr{
				pos: position{line: 238, col: 8, offset: 7036},
				expr: &charClassMatcher{
					pos:        position{line: 238, col: 10, offset: 7038},
					val:        "[a-zA-Z0-9]",
					ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 240, col: 1, offset: 7052},
			expr: &zeroOrMoreExpr{
				pos: position{line: 240, col: 16, offset: 7067},
				expr: &choiceExpr{
					pos: position{line: 240, col: 17, offset: 7068},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 240, col: 17, offset: 7068},
							name: "SPACE",
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 25, offset: 7076},
							name: "COMMENT_TRAIL",
						},
						&ruleRefExpr{
							pos:  position{line: 240, col: 41, offset: 7092},
							name: "COMMENT_BLOCK",
						},
					},
				},
			},
		},
		{
			name: "COMMENT_TRAIL",
			pos:  position{line: 241, col: 1, offset: 7108},
			expr: &seqExpr{
				pos: position{line: 241, col: 19, offset: 7126},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 241, col: 19, offset: 7126},
						val:        "--",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 241, col: 24, offset: 7131},
						expr: &seqExpr{
							pos: position{line: 241, col: 25, offset: 7132},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 241, col: 25, offset: 7132},
									expr: &litMatcher{
										pos:        position{line: 241, col: 26, offset: 7133},
										val:        "\n",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 241, col: 31, offset: 7138,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "COMMENT_BLOCK",
			pos:  position{line: 242, col: 1, offset: 7142},
			expr: &seqExpr{
				pos: position{line: 242, col: 19, offset: 7160},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 242, col: 19, offset: 7160},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 242, col: 24, offset: 7165},
						expr: &seqExpr{
							pos: position{line: 242, col: 25, offset: 7166},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 242, col: 25, offset: 7166},
									expr: &litMatcher{
										pos:        position{line: 242, col: 26, offset: 7167},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&anyMatcher{
									line: 242, col: 31, offset: 7172,
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 242, col: 35, offset: 7176},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "KEY",
			pos:  position{line: 243, col: 1, offset: 7181},
			expr: &notExpr{
				pos: position{line: 243, col: 16, offset: 7196},
				expr: &ruleRefExpr{
					pos:  position{line: 243, col: 17, offset: 7197},
					name: "ID_CONT",
				},
			},
		},
		{
			name: "SPACE",
			pos:  position{line: 244, col: 1, offset: 7205},
			expr: &choiceExpr{
				pos: position{line: 244, col: 16, offset: 7220},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 244, col: 16, offset: 7220},
						val:        " ",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 244, col: 22, offset: 7226},
						val:        "\n",
						ignoreCase: false,
					},
					&litMatcher{
						pos:        position{line: 244, col: 29, offset: 7233},
						val:        "\t",
						ignoreCase: false,
					},
				},
			},
		},
	},
}

func (c *current) onRoot1(expr interface{}) (interface{}, error) {
	return expr, nil
}

func (p *parser) callonRoot1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRoot1(stack["expr"])
}

func (c *current) onSelectStatement1(list, predicateClause, propertyClause interface{}) (interface{}, error) {
	return makeSelect(list, predicateClause, propertyClause) // TODO: makeSelect
}

func (p *parser) callonSelectStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSelectStatement1(stack["list"], stack["predicateClause"], stack["propertyClause"])
}

func (c *current) onDescribeStatement1(statement interface{}) (interface{}, error) {
	return statement, nil
}

func (p *parser) callonDescribeStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDescribeStatement1(stack["statement"])
}

func (c *current) onDescribeAllStatement1(matchClause interface{}) (interface{}, error) {
	return makeDescribeAll(matchClause) // TODO: makeDescribeAll
}

func (p *parser) callonDescribeAllStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDescribeAllStatement1(stack["matchClause"])
}

func (c *current) onEmptyMatchClause1() (interface{}, error) {
	return makeNullMatchClause(), nil
}

func (p *parser) callonEmptyMatchClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyMatchClause1()
}

func (c *current) onMatchClause1(literal interface{}) (interface{}, error) {
	return makeMatchClause(literal) // TODO: makeMatchClause
}

func (p *parser) callonMatchClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMatchClause1(stack["literal"])
}

func (c *current) onDescribeMetrics1(tagName, tagValue interface{}) (interface{}, error) {
	return makeDescribeMetrics(tagName, tagValue) // TODO: makeDescribeMetrics
}

func (p *parser) callonDescribeMetrics1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDescribeMetrics1(stack["tagName"], stack["tagValue"])
}

func (c *current) onDescribeSingleStatement1(name, predicateClause interface{}) (interface{}, error) {
	return makeDescribeSingleStatement(name, predicateClause) // TODO: makeDescribeSingleStatement
}

func (p *parser) callonDescribeSingleStatement1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDescribeSingleStatement1(stack["name"], stack["predicateClause"])
}

func (c *current) onUncheckedPropertyClause4(key, value interface{}) (interface{}, error) {
	return makePropertyKeyValuePair(key, value)
}

func (p *parser) callonUncheckedPropertyClause4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUncheckedPropertyClause4(stack["key"], stack["value"])
}

func (c *current) onUncheckedPropertyClause1(propertyList interface{}) (interface{}, error) {
	return makePropertyClause(propertyList) // TODO: makePropertyClause
}

func (p *parser) callonUncheckedPropertyClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUncheckedPropertyClause1(stack["propertyList"])
}

func (c *current) onPropertyClause1(clause interface{}) (interface{}, error) {
	return checkPropertyClause(clause) // TODO: checkPropertyClause
}

func (p *parser) callonPropertyClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPropertyClause1(stack["clause"])
}

func (c *current) onOptionalPredicateClause3() (interface{}, error) {
	makeNullPredicate() /* TODO: makeNullPredicate */
}

func (p *parser) callonOptionalPredicateClause3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOptionalPredicateClause3()
}

func (c *current) onExpressionList7(expression interface{}) (interface{}, error) {
	return expression, nil
}

func (p *parser) callonExpressionList7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionList7(stack["expression"])
}

func (c *current) onExpressionList1(first, rest interface{}) (interface{}, error) {
	return append([]function.Expression{first}, rest...), nil // TODO: does this work?
}

func (p *parser) callonExpressionList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionList1(stack["first"], stack["rest"])
}

func (c *current) onExpression1(sum, pipe interface{}) (interface{}, error) {
	return addPipe(sum, pipe) /* TODO: implement addPipe */
}

func (p *parser) callonExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpression1(stack["sum"], stack["pipe"])
}

func (c *current) onExpressionSum13(op interface{}) (interface{}, error) {
	return op, nil
}

func (p *parser) callonExpressionSum13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionSum13(stack["op"])
}

func (c *current) onExpressionSum7(pipes, op, right interface{}) (interface{}, error) {
	return makeExpressionSuffix(pipes, op, right) /* TODO: makeExpressionSuffix */
}

func (p *parser) callonExpressionSum7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionSum7(stack["pipes"], stack["op"], stack["right"])
}

func (c *current) onExpressionSum1(first, suffixes interface{}) (interface{}, error) {
	return makeOperator(first, suffixes) /* TODO: makeOperator */
}

func (p *parser) callonExpressionSum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionSum1(stack["first"], stack["suffixes"])
}

func (c *current) onExpressionProduct13(op interface{}) (interface{}, error) {
	return op, nil
}

func (p *parser) callonExpressionProduct13() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionProduct13(stack["op"])
}

func (c *current) onExpressionProduct7(pipes, op, right interface{}) (interface{}, error) {
	return makeExpressionSuffix(pipes, op, right)
}

func (p *parser) callonExpressionProduct7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionProduct7(stack["pipes"], stack["op"], stack["right"])
}

func (c *current) onExpressionProduct1(first, suffixes interface{}) (interface{}, error) {
	return makeOperator(first, suffixes)
}

func (p *parser) callonExpressionProduct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionProduct1(stack["first"], stack["suffixes"])
}

func (c *current) onOnePipe11() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonOnePipe11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOnePipe11()
}

func (c *current) onOnePipe1(name, arguments interface{}) (interface{}, error) {
	return makeOnePipe(name, arguments)
}

func (p *parser) callonOnePipe1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOnePipe1(stack["name"], stack["arguments"])
}

func (c *current) onCallArguments17() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonCallArguments17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArguments17()
}

func (c *current) onCallArguments1(arguments interface{}) (interface{}, error) {
	return arguments, nil
}

func (p *parser) callonCallArguments1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArguments1(stack["arguments"])
}

func (c *current) onExpressionAtom1(core, annotation interface{}) (interface{}, error) {
	if annotation != "" {
		return makeAnnotation(core, annotation), nil
	}
	return core, nil

}

func (p *parser) callonExpressionAtom1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionAtom1(stack["core"], stack["annotation"])
}

func (c *current) onExpressionRaw4(item interface{}) (interface{}, error) {
	return item, nil
}

func (p *parser) callonExpressionRaw4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionRaw4(stack["item"])
}

func (c *current) onExpressionAnnotationRequired1(contents interface{}) (interface{}, error) {
	return contents, nil
}

func (p *parser) callonExpressionAnnotationRequired1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionAnnotationRequired1(stack["contents"])
}

func (c *current) onExpressionAnnotation3() (interface{}, error) {
	return "", nil
}

func (p *parser) callonExpressionAnnotation3() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionAnnotation3()
}

func (c *current) onOptionalGroupBy4() (interface{}, error) {
	return function.Groups{}, nil
}

func (p *parser) callonOptionalGroupBy4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onOptionalGroupBy4()
}

func (c *current) onExpressionFunction1(name, arguments interface{}) (interface{}, error) {
	return makeFunctionCall(name, arguments)
}

func (p *parser) callonExpressionFunction1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionFunction1(stack["name"], stack["arguments"])
}

func (c *current) onExpressionMetric7(predicate interface{}) (interface{}, error) {
	return predicate, nil
}

func (p *parser) callonExpressionMetric7() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionMetric7(stack["predicate"])
}

func (c *current) onExpressionMetric15() (interface{}, error) {
	return makeNullPredicate()
}

func (p *parser) callonExpressionMetric15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionMetric15()
}

func (c *current) onExpressionMetric1(name, predicate interface{}) (interface{}, error) {
	return makeMetricExpression(name, predicate)
}

func (p *parser) callonExpressionMetric1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionMetric1(stack["name"], stack["predicate"])
}

func (c *current) onGroupByClause1(first, rest interface{}) (interface{}, error) {
	return makeGroupBy(first, rest)
}

func (p *parser) callonGroupByClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGroupByClause1(stack["first"], stack["rest"])
}

func (c *current) onCollapseByClause1(first, rest interface{}) (interface{}, error) {
	return makeCollapseBy(first, rest)
}

func (p *parser) callonCollapseByClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCollapseByClause1(stack["first"], stack["rest"])
}

func (c *current) onPredicateClause1(predicate interface{}) (interface{}, error) {
	return predicate, nil
}

func (p *parser) callonPredicateClause1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPredicateClause1(stack["predicate"])
}

func (c *current) onPredicateDisjunction2(left, right interface{}) (interface{}, error) {
	return makeOrPredicate(left, right)
}

func (p *parser) callonPredicateDisjunction2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPredicateDisjunction2(stack["left"], stack["right"])
}

func (c *current) onPredicateConjunction2(left, right interface{}) (interface{}, error) {
	return makeAndPredicate(left, right)
}

func (p *parser) callonPredicateConjunction2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPredicateConjunction2(stack["left"], stack["right"])
}

func (c *current) onPredicateAtom2(atom interface{}) (interface{}, error) {
	return makeNotPredicate(atom)
}

func (p *parser) callonPredicateAtom2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPredicateAtom2(stack["atom"])
}

func (c *current) onPredicateAtom8(predicate interface{}) (interface{}, error) {
	return predicate, nil
}

func (p *parser) callonPredicateAtom8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPredicateAtom8(stack["predicate"])
}

func (c *current) onTagMatcher2(name, literal interface{}) (interface{}, error) {
	return makeLiteralMatcher(name, literal), nil
}

func (p *parser) callonTagMatcher2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagMatcher2(stack["name"], stack["literal"])
}

func (c *current) onTagMatcher10(name, literal interface{}) (interface{}, error) {
	return makeNotPredicate(atom, makeLiteralMatcher(name, literal)), nil
}

func (p *parser) callonTagMatcher10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagMatcher10(stack["name"], stack["literal"])
}

func (c *current) onTagMatcher18(name, literal interface{}) (interface{}, error) {
	return makeRegexMatcher(name, literal), nil
}

func (p *parser) callonTagMatcher18() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagMatcher18(stack["name"], stack["literal"])
}

func (c *current) onTagMatcher27(name, list interface{}) (interface{}, error) {
	return makeListMatcher(name, list), nil
}

func (p *parser) callonTagMatcher27() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagMatcher27(stack["name"], stack["list"])
}

func (c *current) onLiteralString1(contents interface{}) (interface{}, error) {
	return unescape(contents), nil
}

func (p *parser) callonLiteralString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteralString1(stack["contents"])
}

func (c *current) onLiteralList8(literal interface{}) (interface{}, error) {
	return literal, nil
}

func (p *parser) callonLiteralList8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteralList8(stack["literal"])
}

func (c *current) onLiteralList1(first, rest interface{}) (interface{}, error) {
	return makeLiteralList(first, rest), nil
}

func (p *parser) callonLiteralList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteralList1(stack["first"], stack["rest"])
}

func (c *current) onIdentifier2(contents interface{}) (interface{}, error) {
	return contents, nil
}

func (p *parser) callonIdentifier2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier2(stack["contents"])
}

func (c *current) onIdentifier20(segment interface{}) (interface{}, error) {
	return segment, nil
}

func (p *parser) callonIdentifier20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier20(stack["segment"])
}

func (c *current) onIdentifier9(first, rest interface{}) (interface{}, error) {
	return makeIdentifier(first, rest), nil
}

func (p *parser) callonIdentifier9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier9(stack["first"], stack["rest"])
}

func (c *current) onTimestamp2(number, suffix interface{}) (interface{}, error) {
	return number + suffix, nil
}

func (p *parser) callonTimestamp2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestamp2(stack["number"], stack["suffix"])
}

func (c *current) onTimestamp10() (interface{}, error) {
	return "now", nil
}

func (p *parser) callonTimestamp10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTimestamp10()
}

func (c *current) onPropertyKey8() (interface{}, error) {
	return "sample", nil
}

func (p *parser) callonPropertyKey8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPropertyKey8()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	// can't match EOF
	if cur == utf8.RuneError {
		return nil, false
	}
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
